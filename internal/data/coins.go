package data

import (
	"database/sql"
	"time"
)

type CoinModel struct {
	DB *sql.DB
}

type Coin struct {
	CoinID               string    `json:"coin_id"`
	UserID               int64     `json:"user_id"`
	CreatedAt            time.Time `json:"created_at"`
	Symbol               string    `json:"symbol"`
	Amount               float64   `json:"amount"`
	PurchasePriceAverage float64   `json:"purchase_price_average"`
	TotalCost            float64   `json:"total_cost"`
	PNL                  float64   `json:"pnl"`
	Version              int       `json:"version"`
}
