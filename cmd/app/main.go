package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/api"
	"github.com/aalperen0/portfolio-tracker/internal/mail"
	"github.com/aalperen0/portfolio-tracker/internal/model"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

func main() {

	cfg := config.LoadConfig()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

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

	mailer := mail.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender)

	handler := api.NewHandler(*cfg, logger, models, mailer)

	//	mux := http.NewServeMux()
	//mux.HandleFunc("/v1/healthcheck", handler.HealthCheckHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info().Msgf("starting %s server on %s", cfg.Env, srv.Addr)

	err = srv.ListenAndServe()
	logger.Fatal().Err(err)

}
