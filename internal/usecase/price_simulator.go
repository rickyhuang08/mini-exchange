package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rickyhuang08/mini-exchange.git/internal/delivery/websocket"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type PriceSimulator struct {
	hub *websocket.Hub
	stocks []string
	lastPrices map[string]float64
	quit chan struct{}
	Logger *logger.Logger
}

func NewPriceSimulator(hub *websocket.Hub, stocks []string, logger *logger.Logger) *PriceSimulator {
	ps := &PriceSimulator{
		hub: hub,
		stocks: stocks,
		lastPrices: make(map[string]float64),
		quit: make(chan struct{}),
		Logger: logger,
	}
	for _, stock := range stocks {
		ps.lastPrices[stock] = 100.0 + rand.Float64()*50
	}
	return ps
}

func (ps *PriceSimulator) Start() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				ps.updatePrices()
			case <-ps.quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (ps *PriceSimulator) Stop() {
	close(ps.quit)
}

func (ps *PriceSimulator) updatePrices() {
	for _, stock := range ps.stocks {
		change := (rand.Float64() - 0.5) * 2
		ps.lastPrices[stock] += change
		if ps.lastPrices[stock] < 0.01 {
			ps.lastPrices[stock] = 0.01
		}

		ps.hub.BroadcastToStock("ticker", stock, map[string]interface{}{
			"price": ps.lastPrices[stock],
			"time": time.Now().Unix(),
		})
		ps.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Price simulator: updating %s to %f", stock, ps.lastPrices[stock]))
	}
}