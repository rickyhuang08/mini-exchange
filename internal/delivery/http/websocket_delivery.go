package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rickyhuang08/mini-exchange.git/internal/delivery/websocket"
	"github.com/rickyhuang08/mini-exchange.git/pkg/logger"
)

func (h *Handler) WebSocketHandler(c *gin.Context) {
	h.Logger.LogLevel(logger.LogLevelInfo,
		"WebSocket connection attempt from: "+c.Request.RemoteAddr)

	conn, err := websocket.Upgrade(c.Writer, c.Request)
	if err != nil {
		h.Logger.LogLevel(logger.LogLevelError,
			"WebSocket upgrade error: "+err.Error())
		return
	}

	client := &websocket.Client{
		Hub:  h.Hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	h.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}