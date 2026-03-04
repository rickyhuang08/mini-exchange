package domain

import "time"

type Trade struct {
	ID          string    `json:"id"`
	StockCode   string    `json:"stock_code"`
	Price       float64   `json:"price"`
	Quantity    int64     `json:"quantity"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	TradeAt     time.Time `json:"timestamp"`
}
