package config

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// / Initalize database connection, pinging db to check connection
// / is established, if after 5 seconds still not established, cancel
// / process.
// / # Return
// / - db: Returns db driver
// / - error: Returns error if connection couldnt established.

func InitDB(cfg *Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.DB.maxIdleConns)
	db.SetMaxOpenConns(cfg.DB.maxOpenConns)

	maxLifeTime, err := time.ParseDuration(cfg.DB.maxLifeTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(maxLifeTime)

	duration, err := time.ParseDuration(cfg.DB.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// / Initalize redis connection, pinging redis to check connection
// / is established, if after 10 seconds still not established, cancel
// / process.
// / # Return
// / - redis: Returns redis driver
// / - error: Returns error if connection couldnt established.

func InitRedis(cfg *Config, logger zerolog.Logger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to redis")
		return nil, err
	}
	logger.Info().Msg("Connected to redis")

	return rdb, nil
}
