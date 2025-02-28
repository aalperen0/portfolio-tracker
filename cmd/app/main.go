package main

import (
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/api"
	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/mail"
	"github.com/aalperen0/portfolio-tracker/internal/model"
	"github.com/aalperen0/portfolio-tracker/internal/worker"
)

func main() {
	cfg := config.LoadConfig()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	marketData := data.NewCoinClient(cfg.Coins.ApiKey)

	db, err := config.InitDB(cfg)
	if err != nil {
		logger.Fatal().Err(err)
	}
	defer db.Close()
	logger.Info().Msg("database connection pool established")

	rdb, err := config.InitRedis(cfg, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initalize redis")
	}

	models, err := model.NewModels(db, rdb)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initalize models")
	}

	cm := &data.CoinModel{
		DB:  db,
		RDB: rdb,
	}

	pnlUpdater := worker.NewPNLUpdater(cm, marketData, 3*time.Minute, logger)
	pnlUpdater.Start()

	mailer := mail.New(
		cfg.Smtp.Host,
		cfg.Smtp.Port,
		cfg.Smtp.Username,
		cfg.Smtp.Password,
		cfg.Smtp.Sender,
	)

	handler := api.NewHandler(*cfg, logger, models, mailer, marketData)

	err = handler.Serve()
	if err != nil {
		logger.Fatal().Err(err)
	}
}
