package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/api/handlers"
	"github.com/rs/zerolog"
)

func main() {

	cfg := config.LoadConfig()
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	handler := handlers.NewHandler(cfg, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", handler.HealthCheckHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info().Msgf("starting %s server on %s", cfg.Env, srv.Addr)

	err := srv.ListenAndServe()
	logger.Fatal().Err(err)

}
