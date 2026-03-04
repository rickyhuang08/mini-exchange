package entity

type OrderRequest struct {
	StockCode string  `json:"stock_code" binding:"required"`
	Side      string `json:"side" binding:"required,oneof=BUY SELL"`
	Price     float64 `json:"price" binding:"required"`
	Quantity  int64   `json:"quantity" binding:"required"`
}