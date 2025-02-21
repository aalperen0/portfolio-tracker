package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/api"
	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/mail"
	"github.com/aalperen0/portfolio-tracker/internal/model"
)

func main() {
	cfg := config.LoadConfig()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	marketData := data.NewCoinClient(cfg.Coins.ApiKey)

	db, err := config.InitDB(cfg)
	if err != nil {
		logger.Fatal().Err(err)
	}

	models, err := model.NewModels(db)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initalize models")
	}

	defer db.Close()
	logger.Info().Msg("database connection pool established")

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
