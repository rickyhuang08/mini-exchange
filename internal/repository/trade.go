package repository

import (
	"sync"

	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
)

type TradeInterface interface {
	Save(trade *domain.Trade) error
	FindByID(id string) (*domain.Trade, error)
	FindAll(filter map[string]interface{}) ([]*domain.Trade, error)
	FindRecent(stockCode string, limit int) ([]*domain.Trade, error)
	Update(trade *domain.Trade) error
}

type InMemoryTradeRepository struct {
	mu sync.RWMutex
	trades map[string]*domain.Trade
}

func NewInMemoryTradeRepository() TradeInterface {
	return &InMemoryTradeRepository{
		trades: make(map[string]*domain.Trade),
	}
}

func (r *InMemoryTradeRepository) Save(trade *domain.Trade) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trades[trade.ID] = trade
	return nil
}

func (r *InMemoryTradeRepository) FindByID(id string) (*domain.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	trade, exists := r.trades[id]
	if !exists {
		return nil, nil
	}
	return trade, nil
}

func (r *InMemoryTradeRepository) FindAll(filter map[string]interface{}) ([]*domain.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Trade
	for _, trade := range r.trades {
		if stock, ok := filter["stock_code"]; ok && trade.StockCode != stock {
			continue
		}
		result = append(result, trade)
	}
	return result, nil
}

func (r *InMemoryTradeRepository) FindRecent(stockCode string, limit int) ([]*domain.Trade, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var trades []*domain.Trade
	for _, trade := range r.trades {
		if trade.StockCode == stockCode {
			trades = append(trades, trade)
		}
	}
	if len(trades) > limit {
		trades = trades[len(trades)-limit:]
	}
	return trades, nil
}

func (r *InMemoryTradeRepository) Update(trade *domain.Trade) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.trades[trade.ID]; !exists {
		return nil
	}
	r.trades[trade.ID] = trade
	return nil
}