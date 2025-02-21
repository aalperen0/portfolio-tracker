package api

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/aalperen0/portfolio-tracker/config"
	"github.com/aalperen0/portfolio-tracker/internal/data"
	"github.com/aalperen0/portfolio-tracker/internal/mail"
	"github.com/aalperen0/portfolio-tracker/internal/model"
)

type Handler struct {
	config     config.Config
	logger     zerolog.Logger
	models     model.Models
	mailer     mail.Mailer
	wg         sync.WaitGroup
	marketData *data.Client
}

func NewHandler(
	cfg config.Config,
	logger zerolog.Logger,
	models model.Models,
	mailer mail.Mailer,
	marketData *data.Client,
) *Handler {
	return &Handler{
		config:     cfg,
		logger:     logger,
		models:     models,
		mailer:     mailer,
		marketData: marketData,
	}
}
