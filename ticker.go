package yobit

type TickerInfoResponse struct {
	Server_time int               `json:"server_time"`
	Pairs       map[string]Ticker `json:"pairs"`
}

type Ticker struct {
	Decimal_places float64 `json:"decimal_places"`
	Min_price      float64 `json:"min_price"`
	Max_price      float64 `json:"max_price"`
	Min_amount     float64 `json:"min_amount"`
	Min_total      float64 `json:"min_total"`
	Hidden         int     `json:"hidden"`
	Fee            float64 `json:"fee"`
	Fee_buyer      float64 `json:"fee_buyer"`
	Fee_seller     float64 `json:"fee_seller"`
}
