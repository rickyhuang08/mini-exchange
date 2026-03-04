package external

import (
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rickyhuang08/mini-exchange.git/config"
	"github.com/rickyhuang08/mini-exchange.git/internal/domain"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type FinnhubAdapter struct {
	FinnhubCfg config.Finnhub
	Logger *logger.Logger
	TradeOut chan <- *domain.Trade
	Conn *websocket.Conn
	Done chan struct{}
}

func NewFinnhubAdapter(
	finnhubCfg config.Finnhub, 
	logger *logger.Logger,
	tradeOut chan <- *domain.Trade,
) *FinnhubAdapter {
	return &FinnhubAdapter{
		FinnhubCfg: finnhubCfg,
		Logger: logger,
		TradeOut: tradeOut,
		Done: make(chan struct{}),
	}
}

func (fa *FinnhubAdapter) Run() {
	fa.Logger.LogLevel(logger.LogLevelInfo, "Starting Finnhub adapter")
	url := url.URL{Scheme: fa.FinnhubCfg.Scheme, Host: fa.FinnhubCfg.Host, Path: "/", RawQuery: "token=" + fa.FinnhubCfg.ApiKey}
	var err error

	fa.Logger.LogLevel(logger.LogLevelInfo, "Connecting to Finnhub WebSocket at: "+url.String())
	fa.Conn, _, err = websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		fa.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("Finnhub WebSocket dial error: %v", err))
		return
	}
	defer fa.Conn.Close()
	defer close(fa.Done)

	// Subscribe to stock symbols
	fa.Logger.LogLevel(logger.LogLevelInfo, "Subscribing to stock symbols")
	for _, symbol := range fa.FinnhubCfg.Symbols {
		subMsg := map[string]interface{}{
			"type": "subscribe",
			"symbol": symbol,
		}
		if err := fa.Conn.WriteJSON(subMsg); err != nil {
			fa.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("Finnhub subscribe error for symbol %s: %v", symbol, err))
			return
		}
		fa.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Subscribed to Finnhub symbol: %s", symbol))
	}
	
	for {
		select {
			case <-fa.Done:
				return
			default:
				var msg map[string]interface{}
				err := fa.Conn.ReadJSON(&msg)
				if err != nil {
					fa.Logger.LogLevel(logger.LogLevelError, fmt.Sprintf("Finnhub read error: %v", err))
					return
				}

				// 🔍 LOG RAW MESSAGE
        		fa.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("🔍 Raw Finnhub message: %v", msg))

				if msgType, oke := msg["type"]; oke && msgType == "trade" {
					fa.handleTrade(msg)
				}
			
		}
	}
}

func (fa *FinnhubAdapter) handleTrade(raw map[string]interface{}) {
	// 🔍 LOG bahwa handleTrade dipanggil
    fa.Logger.LogLevel(logger.LogLevelInfo, "🔍 handleTrade called with raw: "+fmt.Sprintf("%v", raw))

	data, ok := raw["data"].([]interface{})
	if !ok || len(data) == 0 {
		fa.Logger.LogLevel(logger.LogLevelInfo, "🔍 No trade data or invalid format")
		return
	}

	tradeData, ok := data[0].(map[string]interface{})
	if !ok {
		fa.Logger.LogLevel(logger.LogLevelError, "🔍 Invalid trade data format")
		return
	}

	price, _ := tradeData["p"].(float64)
	volume, _ := tradeData["v"].(float64)
	symbol, _ := tradeData["s"].(string)
	timestamp, _ := tradeData["t"].(float64)

	trade := &domain.Trade{
		ID: uuid.New().String(),
		StockCode: symbol,
		Price: price,
		Quantity: int64(volume),
		TradeAt: time.Unix(0, int64(timestamp)*int64(time.Millisecond)),
	}

	// 🔍 LOG sebelum dikirim
	fa.Logger.LogLevel(logger.LogLevelInfo, "🔍 Forwarding trade: "+fmt.Sprintf("%s %.2f %d", symbol, price, int64(volume)))

	select {
	case fa.TradeOut <- trade:
	default:
		fa.Logger.LogLevel(logger.LogLevelError, "🔍 Finnhub trade channel is full, dropping trade: "+fmt.Sprintf("%v", trade))
	}
}

func (fa *FinnhubAdapter) Stop() {
	close(fa.Done)
	if fa.Conn != nil {
		fa.Conn.Close()
	}
	fa.Logger.LogLevel(logger.LogLevelInfo, "Finnhub adapter stopped")
}