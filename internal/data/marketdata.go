package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

type CoinMarketData struct {
	MarketCapRank     int64   `json:"market_cap_rank"`
	Symbol            string  `json:"symbol"`
	ID                string  `json:"id"`
	CurrentPrice      float64 `json:"current_price"`
	MarketCap         float64 `json:"market_cap"`
	PriceChange24h    float64 `json:"price_change_24h"`
	CirculatingSupply float64 `json:"circulating_supply"`
	MaxSupply         float64 `json:"max_supply"`
	ATH               float64 `json:"ath"`
	LastUpdated       string  `json:"last_updated"`
}

func NewCoinClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://api.coingecko.com/api/v3",
		httpClient: &http.Client{},
	}
}

func (c *Client) GetCoinMarkets(
	currency string,
	filters Filters,
) ([]CoinMarketData, error) {
	query := url.Values{}

	query.Add("vs_currency", currency)

	if filters.Ids != "" {
		query.Add("ids", filters.Ids)
	}

	if filters.Page > 0 {
		query.Add("page", strconv.Itoa(filters.Page))
	}
	if filters.PerPage > 0 {
		query.Add("per_page", strconv.Itoa(filters.PerPage))
	}
	if filters.Order != "" {
		query.Add("order", filters.Order)
	}

	searchUrl := fmt.Sprintf("%s/coins/markets?%s", c.baseURL, query.Encode())

	if len(query) > 0 {
		searchUrl = fmt.Sprintf("%s&%s", searchUrl, query.Encode())
	}

	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("API ERROR (status %d): %s", res.StatusCode, string(body))
	}

	var coins []CoinMarketData
	if err := json.NewDecoder(res.Body).Decode(&coins); err != nil {
		return nil, fmt.Errorf("decoding response %w", err)
	}
	return coins, nil
}

// ////////////////////////////////////////////

func (c *Client) GetCoinCurrentPriceAndSymbol(coinID string) (float64, string, error) {
	url := fmt.Sprintf("%s/coins/%s", c.baseURL, coinID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", c.apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// return 0, "", fmt.Errorf("failed to get coin data: status %d", res.StatusCode)
		return 0, "", validator.ErrRecordNotFound
	}

	var response struct {
		Symbol     string `json:"symbol"`
		MarketData struct {
			CurrentPrice struct {
				USD float64 `json:"usd"`
			} `json:"current_price"`
		} `json:"market_data"`
	}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return 0, "", err
	}
	return response.MarketData.CurrentPrice.USD, response.Symbol, nil
}
