package handlers

import (
	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/rs/zerolog"
)

type Handler struct {
	config config.Config
	logger zerolog.Logger
}

func NewHandler(cfg config.Config, logger zerolog.Logger) *Handler {
	return &Handler{
		config: cfg,
		logger: logger,
	}
}
