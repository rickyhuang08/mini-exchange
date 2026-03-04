package usecase

import (
	"github.com/rickyhuang08/mini-exchange.git/internal/repository"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type MarketUsecase struct {
	bookRepo repository.OrderBookInterface
	tradeRepo repository.TradeInterface
	Logger *logger.Logger
}

type MarketSnapshot struct {
	StockCode string `json:"stock_code"`
	LastPrice float64 `json:"last_price"`
	Change float64 `json:"change"`
	Volume int64 `json:"volume"`
	OrderBook *OrderBookDTO `json:"order_book"`
	RecentTrades []*TradeDTO `json:"recent_trades"`
}

type OrderBookDTO struct {
	Bids [][2]interface{} `json:"bids"`
	Asks [][2]interface{} `json:"asks"`
}

type TradeDTO struct {
	Price float64 `json:"price"`
	Quantity int64 `json:"quantity"`
	Time int64 `json:"time"`
}

func NewMarketUsecase(
	bookRepo repository.OrderBookInterface, 
	tradeRepo repository.TradeInterface,
	logger *logger.Logger,
) *MarketUsecase {
	return &MarketUsecase{
		bookRepo: bookRepo,
		tradeRepo: tradeRepo,
		Logger: logger,
	}
}

func (mu *MarketUsecase) GetSnapshot(stockCode string) (*MarketSnapshot, error) {
	mu.Logger.LogLevel(logger.LogLevelInfo, "GetSnapshot is Running")

	book, err := mu.bookRepo.Get(stockCode)
	mu.Logger.LogLevel(logger.LogLevelInfo, "Retrieved order book for stock: "+stockCode)
	if err != nil {
		mu.Logger.LogLevel(logger.LogLevelError, "Error retrieving order book: "+err.Error())
		return nil, err
	}

	trades, err := mu.tradeRepo.FindRecent(stockCode, 10)
	mu.Logger.LogLevel(logger.LogLevelInfo, "Retrieved recent trades for stock: "+stockCode)
	if err != nil {
		mu.Logger.LogLevel(logger.LogLevelError, "Error retrieving recent trades: "+err.Error())
		return nil, err
	}

	bids := make([][2]interface{}, len(book.Bids))
	for i, order := range book.Bids {
		bids[i] = [2]interface{}{order.Price, order.Quantity}
	}
	mu.Logger.LogLevel(logger.LogLevelInfo, "Processed bids for stock: "+stockCode)

	asks := make([][2]interface{}, len(book.Asks))
	for i, order := range book.Asks {
		asks[i] = [2]interface{}{order.Price, order.Quantity}
	}
	mu.Logger.LogLevel(logger.LogLevelInfo, "Processed asks for stock: "+stockCode)

	recentTrades := make([]*TradeDTO, len(trades))
	for i, trade := range trades {
		recentTrades[i] = &TradeDTO{
			Price: trade.Price,
			Quantity: trade.Quantity,
			Time: trade.TradeAt.Unix(),
		}
	}
	mu.Logger.LogLevel(logger.LogLevelInfo, "Processed recent trades for stock: "+stockCode)

	var lastPrice float64
	if len(trades) > 0 {
		lastPrice = trades[0].Price
	}
	mu.Logger.LogLevel(logger.LogLevelInfo, "Determined last price for stock: "+stockCode)

	var volume int64
	for _, trade := range trades {
		volume += trade.Quantity
	}
	mu.Logger.LogLevel(logger.LogLevelInfo, "Determined volume for stock: "+stockCode)

	change := 0.0

	snapshot := &MarketSnapshot{
		StockCode: stockCode,
		LastPrice: lastPrice,
		Change: change,
		Volume: volume,
		OrderBook: &OrderBookDTO{
			Bids: bids,
			Asks: asks,
		},
		RecentTrades: recentTrades,
	}

	mu.Logger.LogLevel(logger.LogLevelInfo, "Returning market snapshot for stock: "+stockCode)
	return snapshot, nil
}