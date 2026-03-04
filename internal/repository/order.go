package repository

import (
	"sync"

	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
)

type OrderInterface interface {
	Save(order *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	FindAll(filter map[string]interface{}) ([]*domain.Order, error)
	Update(order *domain.Order) error
}

type InMemoryOrderRepository struct {
	mu sync.RWMutex
	orders map[string]*domain.Order
}

func NewInMemoryOrderRepository() OrderInterface {
	return &InMemoryOrderRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *InMemoryOrderRepository) Save(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) FindByID(id string) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, exists := r.orders[id]
	if !exists {
		return nil, nil
	}
	return order, nil
}

func (r *InMemoryOrderRepository) FindAll(filter map[string]interface{}) ([]*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*domain.Order
	for _, order := range r.orders {
		if stock, ok := filter["stock_code"]; ok && order.StockCode != stock {
			continue
		}
		if status, ok := filter["status"]; ok && order.Status != status {
			continue
		}
		result = append(result, order)
	}
	return result, nil
}

func (r *InMemoryOrderRepository) Update(order *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.orders[order.ID]; !exists {
		return nil
	}
	r.orders[order.ID] = order
	return nil
}