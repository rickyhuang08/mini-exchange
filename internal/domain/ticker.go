package domain

import "time"

type Ticker struct {
	StockCode string  `json:"stock_code"`
	LastPrice float64 `json:"last_price"`
	Change    float64 `json:"change"`
	Volume    int64   `json:"volume"`
	UpdatedAt time.Time  `json:"updated_at"`
}
