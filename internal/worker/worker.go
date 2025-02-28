package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/aalperen0/portfolio-tracker/internal/data"
)

type PNLUpdater struct {
	coinModel *data.CoinModel
	client    *data.Client
	interval  time.Duration
	logger    zerolog.Logger
}

func NewPNLUpdater(
	coinModel *data.CoinModel,
	client *data.Client,
	interval time.Duration,
	logger zerolog.Logger,
) *PNLUpdater {
	return &PNLUpdater{
		coinModel: coinModel,
		client:    client,
		interval:  interval,
		logger:    logger,
	}
}

func (p *PNLUpdater) Start() {
	p.logger.Info().Msg("Starting worker...")
	go p.processQueue()

	go p.scheduleUpdates()
}

func (p *PNLUpdater) scheduleUpdates() {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for range ticker.C {
		p.logger.Info().Msg("Enqueuing PNL updates for all coins")
		if err := p.coinModel.EnqueuePNLUpdates(); err != nil {
			p.logger.Err(err).Msgf("Failed to enqueue pnl updates %v", err)
		}
	}
}

func (p *PNLUpdater) processQueue() {
	ctx := context.Background()

	for {

		result, err := p.coinModel.RDB.BLPop(ctx, 5*time.Second, "pnl_queue").Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			p.logger.Err(err).Msgf("Error popping PNL from queue: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(result) < 2 {
			p.logger.Info().Msg("Invalid result from queue")
			continue
		}

		coinID := result[1]
		fmt.Println(coinID)

		currentPrice, _, err := p.client.GetCoinCurrentPriceAndSymbol(coinID)
		if err != nil {
			p.logger.Err(err).Msgf("Error getting current price for %s: %v", coinID, err)
		}

		p.logger.Info().Msgf("Updating PNL for coin %s", coinID)

		if err := p.coinModel.UpdatePNLForCoin(coinID, currentPrice); err != nil {
			p.logger.Err(err).Msgf("Error updating PNL for %s: %v", coinID, err)
		}
	}
}
