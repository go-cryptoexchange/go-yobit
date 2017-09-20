package yobit

type Order struct {
	Price    float64
	Quantity float64
}

type OrderBook struct {
	Asks []Order
	Bids []Order
}

type rawOrderBook struct {
	Asks []interface{} `json:"asks"`
	Bids []interface{} `json:"bids"`
}
