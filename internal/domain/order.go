package domain

import "time"

type Side string
const (
	Buy  Side = "BUY"
	Sell Side = "SELL"
)

type OrderStatus string
const (
	StatusNew       OrderStatus = "NEW"
	StatusPartial   OrderStatus = "PARTIAL_FILL"
	StatusFilled    OrderStatus = "FILLED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID        string      `json:"id"`
	StockCode string      `json:"stock_code"`
	Side      Side        `json:"side"`
	Price     float64     `json:"price"`
	Quantity  int64       `json:"quantity"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}