package api

import (
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
func (h *Handler) AddCoinsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CoinID               string  `json:"coin_id"`
		Amount               float64 `json:"amount"`
		PurchasePriceAverage float64 `json:"purchase_price_average"`
		TotalCost            float64 `json:"total_cost"`
	}

	err := h.readJSON(w, r, &input)
	if err != nil {
		h.badRequestResponse(w, r, err)
		return
	}

	currentPrice, symbol, err := h.marketData.GetCoinCurrentPriceAndSymbol(input.CoinID)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	currentValue := input.Amount * currentPrice
	initialPNL := currentValue - input.TotalCost

	coin := &data.Coin{
		CoinID:               input.CoinID,
		Symbol:               symbol,
		Amount:               input.Amount,
		PurchasePriceAverage: input.PurchasePriceAverage,
		TotalCost:            input.TotalCost,
		PNL:                  initialPNL,
	}

	v := validator.New()
	if data.ValidateCoin(v, coin); !v.Valid() {
		h.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = h.models.Coin.Insert(coin)
	if err != nil {
		h.serverErrorResponse(w, r, err)
		return
	}

	err = h.writeJSON(w, http.StatusCreated, envelope{"coin:": coin}, nil)
	if err != nil {
		h.serverErrorResponse(w, r, err)
	}
}
