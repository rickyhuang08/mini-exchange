package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client

	StockTickerSubs map[string]map[*Client]bool
	TradeSubs       map[string]map[*Client]bool
	OrderBookSubs   map[string]map[*Client]bool
	mu              sync.RWMutex

	Logger *logger.Logger
}

func NewHub(logger *logger.Logger) *Hub {
	return &Hub{
		Clients:         make(map[*Client]bool),
		Broadcast:       make(chan []byte),
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		StockTickerSubs: make(map[string]map[*Client]bool),
		TradeSubs:       make(map[string]map[*Client]bool),
		OrderBookSubs:   make(map[string]map[*Client]bool),
		Logger:          logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Client registered: %p, total clients: %d", client, len(h.Clients)+1))
			h.Clients[client] = true
		case client := <-h.Unregister:
			h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Client Processing unregistered: %p, total clients: %d", client, len(h.Clients)-1))
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)

				h.mu.Lock()
				for _, subs := range h.StockTickerSubs {
					delete(subs, client)
				}
				for _, subs := range h.TradeSubs {
					delete(subs, client)
				}
				for _, subs := range h.OrderBookSubs {
					delete(subs, client)
				}
				h.mu.Unlock()

				log.Printf("Client unregistered: %p, total clients: %d", client, len(h.Clients))
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *Hub) Subscribe(client *Client, channel string, stockCode string) {
	h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Client subscribing to %s:%s", channel, stockCode))
	h.mu.Lock()
	defer h.mu.Unlock()

	var subs map[*Client]bool
	switch channel {
	case "ticker":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Subscribing to ticker for stock: %s", stockCode))
		if h.StockTickerSubs[stockCode] == nil {
			h.StockTickerSubs[stockCode] = make(map[*Client]bool)
		}
		subs = h.StockTickerSubs[stockCode]
	case "trade":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Subscribing to trade for stock: %s", stockCode))
		if h.TradeSubs[stockCode] == nil {
			h.TradeSubs[stockCode] = make(map[*Client]bool)
		}
		subs = h.TradeSubs[stockCode]
	case "orderbook":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Subscribing to orderbook for stock: %s", stockCode))
		if h.OrderBookSubs[stockCode] == nil {
			h.OrderBookSubs[stockCode] = make(map[*Client]bool)
		}
		subs = h.OrderBookSubs[stockCode]
	}
	subs[client] = true
	h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("subscribed client %p to %s:%s, now %d subscribers for this stock", client, channel, stockCode, len(subs)))

}

func (h *Hub) Unsubscribe(client *Client, channel string, stockCode string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Client Unsubscribe to %s:%s", channel, stockCode))

	switch channel {
	case "ticker":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Unsubscribing from ticker for stock: %s", stockCode))
		if subs, ok := h.StockTickerSubs[stockCode]; ok {
			delete(subs, client)
		}
	case "trade":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Unsubscribing from trade for stock: %s", stockCode))
		if subs, ok := h.TradeSubs[stockCode]; ok {
			delete(subs, client)
		}
	case "orderbook":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Unsubscribing from orderbook for stock: %s", stockCode))
		if subs, ok := h.OrderBookSubs[stockCode]; ok {
			delete(subs, client)
		}
	}
}

func (h *Hub) BroadcastToStock(channel string, stockCode string, message interface{}) {
	h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Broadcasting to stock %s:%s", channel, stockCode))
	data, _ := json.Marshal(map[string]interface{}{
		"channel": channel,
		"stock":   stockCode,
		"data":    message,
	})
	h.mu.RLock()
	defer h.mu.RUnlock()
	var subs map[*Client]bool
	switch channel {
	case "ticker":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Broadcasting to ticker for stock: %s", stockCode))
		subs = h.StockTickerSubs[stockCode]
	case "trade":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Broadcasting to trade for stock: %s", stockCode))
		subs = h.TradeSubs[stockCode]
	case "orderbook":
		h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Broadcasting to orderbook for stock: %s", stockCode))
		subs = h.OrderBookSubs[stockCode]
	}

	h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Broadcasting to stock %s:%s, subscribers count: %d", channel, stockCode, len(subs)))
	
	for client := range subs {
		select {
		case client.Send <- data:
		default:
			h.Logger.LogLevel(logger.LogLevelInfo, fmt.Sprintf("Client %p send buffer full, dropping message", client))
			// close(client.Send)
			// delete(h.Clients, client)
		}
	}
}