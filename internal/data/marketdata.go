package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/aalperen0/portfolio-tracker/internal/cache"
	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	cache      *cache.Cache
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

func NewCoinClient(apiKey string, cache *cache.Cache) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://api.coingecko.com/api/v3",
		httpClient: &http.Client{},
		cache:      cache,
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

// / If a coin stored in the cache, we retrieve coin price and coin symbol.
// / Otherwise retrieve a coin from the CoinGecko api. If user retrieve a coin
// / from the api, the function put the values to the cache for a 5 minute.
// # Parameters
// - coinID (string)
// # Return
// - price (float64)
// - symbol(string)
// - error(record not found)

func (c *Client) GetCoinCurrentPriceAndSymbol(coinID string) (float64, string, error) {
	ctx := context.Background()

	if c.cache != nil {
		cacheKey := "coin:price:" + coinID

		type CachedData struct {
			Price  float64 `json:"price"`
			Symbol string  `json:"symbol"`
		}

		var cachedData CachedData
		found, err := c.cache.Get(ctx, cacheKey, &cachedData)
		if err == nil && found {
			return cachedData.Price, cachedData.Symbol, nil
		}

	}

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

	price := response.MarketData.CurrentPrice.USD
	symbol := response.Symbol

	if c.cache != nil {
		cachedData := struct {
			Price  float64 `json:"price"`
			Symbol string  `json:"symbol"`
		}{
			Price:  price,
			Symbol: symbol,
		}
		err := c.cache.Set(ctx, "coin:price:"+coinID, cachedData, 5*time.Minute)
		if err != nil {
			return 0, "", err
		}
	}

	return price, symbol, nil
}
