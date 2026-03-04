package usecase

import (
	"fmt"

	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
	"github.com/rickyhuang08/mini-exchange.git/internal/repository"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type OrderUsecase struct {
	orderRepo repository.OrderInterface
	tradeRepo repository.TradeInterface
	MatchingEngine *MatchingEngine
	Logger *logger.Logger
}

func NewOrderUsecase(
	orderRepo repository.OrderInterface, 
	tradeRepo repository.TradeInterface,
	me *MatchingEngine,
	logger *logger.Logger,
) *OrderUsecase {
	return &OrderUsecase{
		orderRepo: orderRepo,
		tradeRepo: tradeRepo,
		MatchingEngine: me,
		Logger: logger,
	}
}

func (ou *OrderUsecase) PlaceOrder(order *domain.Order) error {
	ou.Logger.LogLevel(logger.LogLevelInfo, "Placing order for stock: "+order.StockCode)
	if err := ou.orderRepo.Save(order); err != nil {
		ou.Logger.LogLevel(logger.LogLevelError, "Error saving order: "+err.Error())
		return err
	}
	ou.Logger.LogLevel(logger.LogLevelInfo, "Order saved successfully for stock: "+order.StockCode)
	return ou.MatchingEngine.SubmitOrder(order)
}

func (ou *OrderUsecase) ListOrders(filter map[string]interface{}) ([]*domain.Order, error) {
	ou.Logger.LogLevel(logger.LogLevelInfo, "Listing orders with filter: "+fmt.Sprintf("%v", filter))
	return ou.orderRepo.FindAll(filter)
}

func (ou *OrderUsecase) GetTradeHistory(stockCode string) ([]*domain.Trade, error) {
	ou.Logger.LogLevel(logger.LogLevelInfo, "Getting trade history for stock: "+stockCode)
	return ou.tradeRepo.FindAll(map[string]interface{}{"stock_code": stockCode})
}