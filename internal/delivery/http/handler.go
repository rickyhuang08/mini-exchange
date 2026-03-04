package http

import (
	"github.com/rickyhuang08/mini-exchange.git/internal/delivery/websocket"
	"github.com/rickyhuang08/mini-exchange.git/internal/usecase"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type Handler struct {
	AuthUsecase   *usecase.AuthUsecase
	OrderUsecase  *usecase.OrderUsecase
	MarketUsecase *usecase.MarketUsecase
	Hub           *websocket.Hub
	Logger        *logger.Logger
}

func NewHandler(
	authUsecase *usecase.AuthUsecase,
	orderUsecase *usecase.OrderUsecase,
	marketUsecase *usecase.MarketUsecase,
	hub *websocket.Hub,
	logger *logger.Logger,
) *Handler {
	return &Handler{
		AuthUsecase:   authUsecase,
		OrderUsecase:  orderUsecase,
		MarketUsecase: marketUsecase,
		Hub:           hub,
		Logger:        logger,
	}
}
