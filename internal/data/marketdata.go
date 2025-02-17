package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (c *Client) GetCoinMarkets(currency string) ([]CoinMarketData, error) {
	url := fmt.Sprintf("%s/coins/markets?vs_currency=%s", c.baseURL, currency)

	req, err := http.NewRequest("GET", url, nil)
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
