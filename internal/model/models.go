package model

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/aalperen0/portfolio-tracker/internal/cache"
	"github.com/aalperen0/portfolio-tracker/internal/data"
)

type Models struct {
	User  data.UserModel
	Token data.TokenModel
	Coin  data.CoinModel
	RDB   *redis.Client
	Cache *cache.Cache
}

func NewModels(
	db *sql.DB,
	rdb *redis.Client,
	cache *cache.Cache,
	logger zerolog.Logger,
) (Models, error) {
	return Models{
		User:  data.UserModel{DB: db},
		Token: data.TokenModel{DB: db},
		Coin:  data.CoinModel{DB: db, Cache: cache, Logger: logger},
		RDB:   rdb,
		Cache: cache,
	}, nil
}
