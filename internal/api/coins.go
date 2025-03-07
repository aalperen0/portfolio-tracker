package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

func (h *Handler) GetCoinsFromMarketHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Currency string `json:"currency"`
		data.Filters
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	qs := r.URL.Query()
	input.Page = h.readURLint(qs, "page", 1, v)
	input.PerPage = h.readURLint(qs, "per_page", 20, v)
	input.Ids = h.readURLstring(qs, "ids", "")
	input.Order = h.readURLstring(qs, "order", "market_cap_desc")

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	coins, err := h.marketData.GetCoinMarkets(input.Currency, input.Filters)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "invalid vs_currency"):
			h.badRequestResponse(w, r, validator.ErrInvalidCurrency)
			return
		default:
			h.serverErrorResponse(w, r, err)
		}
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"coins": coins}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// POST /v1/users/coins
// / Adding coins to portfolio
// / We're asking user to coin name, amount of coin, purchase price
// / calculating initial pnl and total cost
// / fetching the coin from coingecko api, if the coin exists, it adds to db
// / successfully, otherwise we return coin couldn't be found

func (h *Handler) AddCoinsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CoinID        string  `json:"coin_id"`
		Amount        float64 `json:"amount"`
		PurchasePrice float64 `json:"purchase_price"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	user := data.ContextGetUser(r)
	if user.IsAnonymous() {
		h.authenticationRequiredResponse(w, r)
		return
	}

	existingCoin, err := h.models.Coin.GetCoinForUser(input.CoinID, user.ID)
	if err == nil {
		err := h.writeJSON(
			w,
			http.StatusConflict,
			envelope{"error": fmt.Sprintf("%s already exists in portfolio", existingCoin.CoinID)},
			nil,
		)
		if err != nil {
			h.serverErrorResponse(w, r, err)
		}
		return
	} else if !errors.Is(err, validator.ErrRecordNotFound) {
		h.serverErrorResponse(w, r, err)
		return
	}

	currentPrice, symbol, err := h.marketData.GetCoinCurrentPriceAndSymbol(input.CoinID)
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrRecordNotFound):
			h.notFoundResponse(w, r)
			return
		default:
			h.serverErrorResponse(w, r, err)
		}
	}

	// INITIAL Total Cost
	totalCost := input.Amount * input.PurchasePrice

	// INITIAL PNL
	currentValue := input.Amount * currentPrice
	initialPNL := currentValue - totalCost

	coin := &data.Coin{
		CoinID:               input.CoinID,
		UserID:               user.ID,
		Symbol:               symbol,
		Amount:               input.Amount,
		PurchasePriceAverage: input.PurchasePrice,
		TotalCost:            totalCost,
		PNL:                  initialPNL,
	}

	v := validator.New()
	if data.ValidateCoin(v, coin); !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = h.models.Coin.InsertCoin(coin)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, envelope{"coin:": coin}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// / GET /v1/users/coins/:id

func (h *Handler) GetCoinFromPortfolioHandler(w http.ResponseWriter, r *http.Request) {
	coinID, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	user := data.ContextGetUser(r)

	coin, err := h.models.Coin.GetCoinForUser(coinID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrRecordNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return

	}

	err = h.writeJSON(w, http.StatusOK, envelope{"coin": coin}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

func (h *Handler) DeleteCoinFromPortfolioHandler(w http.ResponseWriter, r *http.Request) {
	coinID, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	user := data.ContextGetUser(r)

	err = h.models.Coin.DeleteCoin(coinID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrRecordNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return

	}

	err = h.writeJSON(
		w,
		http.StatusOK,
		envelope{"coin": fmt.Sprintf("%s deleted from portoflio", coinID)},
		nil,
	)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}

// / GET /v1/users/coins
func (h *Handler) GetAllCoinsFromPortfolioHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CoinID string
		data.Filters
	}

	user := data.ContextGetUser(r)

	v := validator.New()

	qs := r.URL.Query()
	input.CoinID = h.readURLstring(qs, "coin", "")
	input.Page = h.readURLint(qs, "page", 1, v)
	input.PerPage = h.readURLint(qs, "per_page", 20, v)
	input.Sort = h.readURLstring(qs, "sort", "amount_asc")
	input.SortList = []string{
		"amount_asc",
		"amount_desc",
		"pnl_asc",
		"pnl_desc",
		"coin_id_asc",
		"coin_id_desc",
	}

	if data.ValidateOtherFilters(v, input.Filters); !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	coins, err := h.models.Coin.GetAllCoinsForUser(input.CoinID, user.ID, input.Filters)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"coins": coins}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}
}

// / UPDATE /v1/users/coins/:id

func (h *Handler) UpdateCoinsHandler(w http.ResponseWriter, r *http.Request) {
	coinID, err := h.readIDParam(r)
	if err != nil {
		h.notFoundResponse(w, r)
		return
	}

	user := data.ContextGetUser(r)

	coin, err := h.models.Coin.GetCoinForUser(coinID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrRecordNotFound):
			h.notFoundResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
	}

	var input struct {
		Amount        float64 `json:"amount"`
		PurchasePrice float64 `json:"purchase_price"`
	}

	err = h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateCoinUpdate(v, input.Amount, input.PurchasePrice); !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	currentPrice, _, err := h.marketData.GetCoinCurrentPriceAndSymbol(coinID)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	calculateFields(coin, input.Amount, input.PurchasePrice)
	calculatePNL(coin, currentPrice)

	err = h.models.Coin.UpdateCoinsForUser(coin)
	if err != nil {
		switch {
		case errors.Is(err, validator.ErrEditConflict):
			h.editConflictResponse(w, r)
		default:
			h.serverErrorResponse(w, r, err)
		}
		return
	}

	err = h.writeJSON(w, http.StatusOK, envelope{"coin": coin}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
