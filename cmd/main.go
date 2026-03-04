package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rickyhuang08/mini-exchange.git/config"
	"github.com/rickyhuang08/mini-exchange.git/internal/adapter/external"
	http_delivery "github.com/rickyhuang08/mini-exchange.git/internal/delivery/http"
	"github.com/rickyhuang08/mini-exchange.git/internal/delivery/websocket"
	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
	"github.com/rickyhuang08/mini-exchange.git/internal/repository"
	"github.com/rickyhuang08/mini-exchange.git/internal/usecase"
	"github.com/rickyhuang08/mini-exchange.git/middleware"
	"github.com/rickyhuang08/mini-exchange.git/pkg/jwt"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)


func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	fmt.Print("config :", cfg)

	// Initialize logger
	loggerM, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer loggerM.AccessFile.Close()
	defer loggerM.ErrorFile.Close()

	// Initialize repositories
	userRepo := repository.NewInMemoryUserRepository()
	orderRepo := repository.NewInMemoryOrderRepository()
	tradeRepo := repository.NewInMemoryTradeRepository()
	bookRepo := repository.NewInMemoryOrderBookRepository()

	// Initialize matching engine
	me := usecase.NewMatchingEngine(orderRepo, bookRepo, tradeRepo, loggerM)

	// Initialize JWT Helpers
	jwtHelper := jwt.NewJWTHelper(loggerM, cfg.JWT.PrivateKeyPath, cfg.JWT.Expiration)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtHelper, loggerM, cfg.JWT.PublicKeyPath)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, tradeRepo, me, loggerM)
	marketUsecase := usecase.NewMarketUsecase(bookRepo, tradeRepo, loggerM)

	// WebSocket hub
	hub := websocket.NewHub(loggerM)
	go hub.Run()

	// Start price simulator
	stocks := cfg.Finnhub.Symbols
	priceSimulator := usecase.NewPriceSimulator(hub, stocks, loggerM)
	// Start the price simulator in a separate goroutine (uncommand below to enable it)
	// priceSimulator.Start()
	defer priceSimulator.Stop()

	// Start Finnhub adapter
	externalTradeChan := make(chan *domain.Trade, 100)
	finndhubAdapter := external.NewFinnhubAdapter(
		cfg.Finnhub,
		loggerM,
		externalTradeChan,
	)
	go finndhubAdapter.Run()
	defer finndhubAdapter.Stop()

	// Broadcast Finnhub trades to WebSocket clients
	go func() {
		for trade := range externalTradeChan {
			_ = me.TradeRepository.Save(trade)

			hub.BroadcastToStock("ticker", trade.StockCode, map[string]interface{}{
				"price": trade.Price,
				"time": trade.TradeAt.Unix(),
			})

			hub.BroadcastToStock("trade", trade.StockCode, trade)
		}	
	}()

	// Broadcast internal trades to WebSocket clients
	go func() {
		for trade := range me.TradeBroadcast {
			hub.BroadcastToStock("trade", trade.StockCode, trade)
		}
	}()

	// Initialize Middleware
	mw := middleware.NewMiddlewareModule(cfg.JWT.PublicKeyPath, loggerM)

	// Start HTTP server
	router := gin.Default()
	handler := http_delivery.NewHandler(authUsecase, orderUsecase, marketUsecase, hub, loggerM)
	http_delivery.RegisterRoutes(router, handler, mw)

	loggerM.LogLevel(logger.LogLevelInfo, "Server started on :8080")
	if err := router.Run(":8080"); err != nil {
		loggerM.LogLevel(logger.LogLevelError, "Server error: "+err.Error())
	}	
}