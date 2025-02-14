package api

import (
	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/mail"
	"github.com/aalperen0/portfolio-tracker/internal/model"
	"github.com/rs/zerolog"
)

type Handler struct {
	config config.Config
	logger zerolog.Logger
	models model.Models
	mailer mail.Mailer
}

func NewHandler(cfg config.Config, logger zerolog.Logger, models model.Models, mailer mail.Mailer) *Handler {
	return &Handler{
		config: cfg,
		logger: logger,
		models: models,
		mailer: mailer,
	}
}
