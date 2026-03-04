package domain

type OrderBook struct {
	StockCode string   `json:"stock_code"`
	Bids      []*Order `json:"bids"`
	Asks      []*Order `json:"asks"`
}
