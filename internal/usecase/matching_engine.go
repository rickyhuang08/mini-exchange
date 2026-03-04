package usecase

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
	"github.com/rickyhuang08/mini-exchange.git/internal/repository"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type MatchingEngine struct {
	OrderRepository     repository.OrderInterface
	OrderBookRepository repository.OrderBookInterface
	TradeRepository     repository.TradeInterface
	StockChans          map[string]chan *domain.Order
	StockMu             sync.RWMutex
	TradeBroadcast      chan *domain.Trade
	Logger              *logger.Logger
}

func NewMatchingEngine(
	orderRepo repository.OrderInterface,
	orderBookRepo repository.OrderBookInterface,
	tradeRepo repository.TradeInterface,
	logger *logger.Logger,
) *MatchingEngine {
	me := &MatchingEngine{
		OrderRepository:     orderRepo,
		OrderBookRepository: orderBookRepo,
		TradeRepository:     tradeRepo,
		StockChans:          make(map[string]chan *domain.Order),
		TradeBroadcast:      make(chan *domain.Trade, 100),
		Logger:              logger,
	}
	return me
}

// StartStockProcessor initializes a goroutine to process orders for a specific stock
func (me *MatchingEngine) StartStockProcessor(stockCode string) {
	me.Logger.LogLevel(logger.LogLevelInfo, "Starting stock processor for stock: "+stockCode)

	me.StockMu.Lock()
	defer me.StockMu.Unlock()
	me.Logger.LogLevel(logger.LogLevelInfo, "Lock acquired for stock: "+stockCode)
	if _, exists := me.StockChans[stockCode]; exists {
		me.Logger.LogLevel(logger.LogLevelInfo, "Stock processor already running for stock: "+stockCode)
		return
	}
	me.Logger.LogLevel(logger.LogLevelInfo, "Creating new stock processor for stock: "+stockCode)
	ch := make(chan *domain.Order, 100)
	me.StockChans[stockCode] = ch

	me.Logger.LogLevel(logger.LogLevelInfo, "New stock processor created for stock: "+stockCode)
	go me.prcoessOrders(stockCode, ch)
}

// SubmitOrder sends a new order to the appropriate stock processor channel
func (me *MatchingEngine) SubmitOrder(order *domain.Order) error {
	me.Logger.LogLevel(logger.LogLevelInfo, "Submitting order for stock: "+order.StockCode)
	me.StartStockProcessor(order.StockCode)
	me.StockMu.RLock()
	ch, _ := me.StockChans[order.StockCode]
	me.StockMu.RUnlock()

	me.Logger.LogLevel(logger.LogLevelInfo, "Sending order to stock processor for stock: "+order.StockCode)
	ch <- order
	return nil
}

// prcoessOrders is the main loop that processes incoming orders for a specific stock
func (me *MatchingEngine) prcoessOrders(stockCode string, orders <-chan *domain.Order) {
	me.Logger.LogLevel(logger.LogLevelInfo, "Processing orders for stock: "+stockCode)
	book, _ := me.OrderBookRepository.Get(stockCode)
	var bids, asks []*domain.Order

	me.Logger.LogLevel(logger.LogLevelInfo, "Starting order processing for stock: "+stockCode)
	for order := range orders {
		me.Logger.LogLevel(logger.LogLevelInfo, "Received order for stock: "+stockCode+" Order ID: "+order.ID)
		if order.Side == domain.Buy {
			bids = append(bids, order)
			sort.Slice(bids, func(i, j int) bool {
				return bids[i].Price > bids[j].Price
			})
		} else {
			asks = append(asks, order)
			sort.Slice(asks, func(i, j int) bool {
				return asks[i].Price < asks[j].Price
			})
		}

		me.Logger.LogLevel(logger.LogLevelInfo, "Before Matching Removing filled orders for stock: "+stockCode)
		bids = me.removeFilled(bids)
		asks = me.removeFilled(asks)

		me.Logger.LogLevel(logger.LogLevelInfo, "Matching orders for stock: "+stockCode)
		trades := me.match(bids, asks)
		for _, trade := range trades {
			me.TradeRepository.Save(trade)
			me.TradeBroadcast <- trade

			me.Logger.LogLevel(logger.LogLevelInfo, "Processing trade for stock: "+stockCode+" Trade ID: "+trade.ID)
			buyOrder, _ := me.OrderRepository.FindByID(trade.BuyOrderID)
			if buyOrder != nil {
				// Ensure trade quantity does not exceed buy order quantity
				if trade.Quantity > buyOrder.Quantity {
					me.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("ERROR: trade quantity %d exceeds buy order %s quantity %d", trade.Quantity, buyOrder.ID, buyOrder.Quantity))
					trade.Quantity = buyOrder.Quantity
				}

				buyOrder.Quantity -= trade.Quantity
				if buyOrder.Quantity <= 0 {
					buyOrder.Status = domain.StatusFilled
				} else {
					buyOrder.Status = domain.StatusPartial
				}
				buyOrder.UpdatedAt = time.Now()
				me.OrderRepository.Update(buyOrder)
			}

			me.Logger.LogLevel(logger.LogLevelInfo, "Updated buy order for stock: "+stockCode+" Order ID: "+trade.BuyOrderID)
			sellOrder, _ := me.OrderRepository.FindByID(trade.SellOrderID)
			if sellOrder != nil {
				sellOrder.Quantity -= trade.Quantity
				if sellOrder.Quantity <= 0 {
					sellOrder.Status = domain.StatusFilled
				} else {
					sellOrder.Status = domain.StatusPartial
				}
				sellOrder.UpdatedAt = time.Now()
				me.OrderRepository.Update(sellOrder)
			}
		}

		me.Logger.LogLevel(logger.LogLevelInfo, "Removing filled orders for stock: "+stockCode)
		bids = me.removeFilled(bids)
		asks = me.removeFilled(asks)

		me.Logger.LogLevel(logger.LogLevelInfo, "Updating order book for stock: "+stockCode)
		book.Bids = bids
		book.Asks = asks
		me.OrderBookRepository.Update(book)
	}
	me.Logger.LogLevel(logger.LogLevelInfo, "Finished processing orders for stock: "+stockCode)
}

// match performs the actual matching logic between buy and sell orders
func (me *MatchingEngine) match(bids, asks []*domain.Order) []*domain.Trade {
	var trades []*domain.Trade
	i, j := 0, 0
	for i < len(bids) && j < len(asks) {
		buy := bids[i]
		sell := asks[j]

		// Skip orders that are already filled (quantity <= 0)
        if buy.Quantity <= 0 {
            i++
            continue
        }
        if sell.Quantity <= 0 {
            j++
            continue
        }

		if buy.Price >= sell.Price {
			quantity := min(buy.Quantity, sell.Quantity)
			trade := &domain.Trade{
				ID:          uuid.New().String(),
				StockCode:   buy.StockCode,
				BuyOrderID:  buy.ID,
				SellOrderID: sell.ID,
				Price:       sell.Price,
				Quantity:    quantity,
				TradeAt:     time.Now(),
			}
			trades = append(trades, trade)

			buy.Quantity -= quantity
			sell.Quantity -= quantity

			if buy.Quantity == 0 {
				i++
			}
			if sell.Quantity == 0 {
				j++
			}
		} else {
			break
		}
	}
	return trades
}

// removeFilled removes orders that have been completely filled (quantity <= 0)
func (me *MatchingEngine) removeFilled(orders []*domain.Order) []*domain.Order {
	var result []*domain.Order
	for _, order := range orders {
		if order.Quantity > 0 {
			result = append(result, order)
		}
	}
	return result
}

// min returns the smaller of two int64 values
func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
