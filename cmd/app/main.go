package main

import (
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/api"
	"github.com/aalperen0/portfolio-tracker/internal/cache"
	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/mail"
	"github.com/aalperen0/portfolio-tracker/internal/model"
	"github.com/aalperen0/portfolio-tracker/internal/worker"
)

func main() {
	///////////////////////////////////////////////////////////////
	// Config initialization
	cfg := config.LoadConfig()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	///////////////////////////////////////////////////////////////
	// Database connection
	db, err := config.InitDB(cfg)
	if err != nil {
		logger.Fatal().Err(err)
	}
	defer db.Close()
	logger.Info().Msg("database connection pool established")

	///////////////////////////////////////////////////////////////
	// Redis initialization
	rdb, err := config.InitRedis(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initalize redis")
	}

	cache := cache.NewCache(rdb, 15*time.Minute)
	marketData := data.NewCoinClient(cfg.Coins.ApiKey, cache)

	///////////////////////////////////////////////////////////////
	// Data models and worker initialization
	models, err := model.NewModels(db, rdb, cache, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initalize models")
	}

	cm := &data.CoinModel{
		DB:  db,
		RDB: rdb,
	}

	pnlUpdater := worker.NewPNLUpdater(cm, marketData, 10*time.Minute, logger)
	pnlUpdater.Start()

	///////////////////////////////////////////////////////////////
	/// Mailer initialization
	mailer := mail.New(
		cfg.Smtp.Host,
		cfg.Smtp.Port,
		cfg.Smtp.Username,
		cfg.Smtp.Password,
		cfg.Smtp.Sender,
	)

	///////////////////////////////////////////////////////////////
	// Server initialization
	handler := api.NewHandler(*cfg, logger, models, mailer, marketData)

	err = handler.Serve()
	if err != nil {
		logger.Fatal().Err(err)
	}
}
