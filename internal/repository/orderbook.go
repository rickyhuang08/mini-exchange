package repository

import (
	"sync"

	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
)

type OrderBookInterface interface {
	Get(StockCode string) (*domain.OrderBook, error)
	Update(book *domain.OrderBook) error
}

type InMemoryOrderBookRepository struct {
	mu sync.RWMutex
	books map[string]*domain.OrderBook
}

func NewInMemoryOrderBookRepository() OrderBookInterface {
	return &InMemoryOrderBookRepository{
		books: make(map[string]*domain.OrderBook),
	}
}

func (r *InMemoryOrderBookRepository) Get(stockCode string) (*domain.OrderBook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	book, exists := r.books[stockCode]
	if !exists {
		book = &domain.OrderBook{
			StockCode: stockCode,
			Bids:      []*domain.Order{},
			Asks:      []*domain.Order{},
		}
		r.books[stockCode] = book	
	}
	return book, nil
}

func (r *InMemoryOrderBookRepository) Update(book *domain.OrderBook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.books[book.StockCode] = book
	return nil
}