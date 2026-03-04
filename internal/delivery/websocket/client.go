package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		log.Printf("Received raw message: %s", message)

		var req struct {
			Action	string `json:"action"`
			Channel string `json:"channel"`
			Stock string `json:"stock"`
		}
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("Invalid JSON: %s", message)
			continue
		}
		log.Printf("Received action: %s, channel: %s, stock: %s", req.Action, req.Channel, req.Stock)
		switch req.Action {
		case "subscribe":
			c.Hub.Subscribe(c, req.Channel, req.Stock)
		case "unsubscribe":
			c.Hub.Unsubscribe(c, req.Channel, req.Stock)
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	log.Printf("WritePump started for client %p", c)
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				log.Printf("WritePump: send channel closed for client %p", c)
				return
			}
			log.Printf("WritePump: sending %d bytes to client %p", len(message), c)
			err := c.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("WritePump error: %v", err)
				return
			}
		}
	}
}