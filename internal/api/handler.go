package api

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/mail"
	"github.com/aalperen0/portfolio-tracker/internal/model"
)

type Handler struct {
	config config.Config
	logger zerolog.Logger
	models model.Models
	mailer mail.Mailer
	wg     sync.WaitGroup
}

func NewHandler(
	cfg config.Config,
	logger zerolog.Logger,
	models model.Models,
	mailer mail.Mailer,
) *Handler {
	return &Handler{
		config: cfg,
		logger: logger,
		models: models,
		mailer: mailer,
	}
}
