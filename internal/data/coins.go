package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/aalperen0/portfolio-tracker/internal/validator"
)

type CoinModel struct {
	DB *sql.DB
}

type Coin struct {
	CoinID               string    `json:"coin_id"`
	UserID               int64     `json:"-"`
	CreatedAt            time.Time `json:"-"`
	Symbol               string    `json:"symbol"`
	Amount               float64   `json:"amount"`
	PurchasePriceAverage float64   `json:"purchase_price_average"`
	TotalCost            float64   `json:"total_cost"`
	PNL                  float64   `json:"pnl"`
	Version              int       `json:"version"`
}

func ValidateCoin(v *validator.Validator, coin *Coin) {
	v.Check(coin.Amount > 0, "amount", "must be greater than zero")
	v.Check(coin.PurchasePriceAverage > 0, "purchase_price_average", "must be greater than zero")
	v.Check(coin.TotalCost > 0, "total_cost", "must be greater than zero")
	v.Check(coin.CoinID != "", "coin_id", "must be provided")
	v.Check(len(coin.CoinID) <= 100, "", "must be not longer than 100 bytes")
}

func (m CoinModel) InsertCoin(coin *Coin) error {
	query := `INSERT INTO coins(coin_id, user_id, symbol, amount, purchase_price_average, total_cost, pnl)
              VALUES($1, $2, $3, $4, $5 ,$6, $7)
              RETURNING created_at, version`

	args := []any{
		coin.CoinID,
		coin.UserID,
		coin.Symbol,
		coin.Amount,
		coin.PurchasePriceAverage,
		coin.TotalCost,
		coin.PNL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&coin.CreatedAt, &coin.Version)
}

// / Get coin from the porfolio according to id of coin
// / Coin id must be string
func (m CoinModel) GetCoinForUser(coinId string, userID int64) (*Coin, error) {
	if coinId == "" {
		return nil, validator.ErrRecordNotFound
	}

	query := `SELECT coin_id, user_id, symbol, amount, purchase_price_average, total_cost, pnl
              FROM coins
              WHERE coin_id = $1 AND user_id = $2`

	var coin Coin

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, coinId, userID).Scan(
		&coin.CoinID,
		&coin.UserID,
		&coin.Symbol,
		&coin.Amount,
		&coin.PurchasePriceAverage,
		&coin.TotalCost,
		&coin.PNL,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, validator.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &coin, nil
}

// Update  coin in the porfolio
func (m CoinModel) UpdateCoin(coin *Coin) error {
	return nil
}

// Delete  coin in the porfolio
func (m CoinModel) DeleteCoin(coinID string, userID int64) error {
	if coinID == "" {
		return validator.ErrRecordNotFound
	}

	query := `DELETE FROM coins WHERE coin_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, coinID, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return validator.ErrRecordNotFound
	}

	return nil
}

func (m CoinModel) GetAllCoinsForUser(
	coinID string,
	userID int64,
	filters Filters,
) ([]*Coin, error) {
	query := fmt.Sprintf(`SELECT coin_id, symbol, amount, purchase_price_average, total_cost, pnl
              FROM coins
              WHERE (coin_id ILIKE $1 OR symbol ILIKE $1 or $1 = '') AND user_id = $2
              ORDER BY %s %s 
              LIMIT $3 OFFSET $4`, filters.SortColumn(), filters.SortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{coinID, userID, filters.Limit(), filters.Offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	coins := []*Coin{}

	for rows.Next() {
		var coin Coin
		err := rows.Scan(
			&coin.CoinID,
			&coin.Symbol,
			&coin.Amount,
			&coin.PurchasePriceAverage,
			&coin.TotalCost,
			&coin.PNL,
		)
		if err != nil {
			return nil, err
		}
		coins = append(coins, &coin)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return coins, nil
}
