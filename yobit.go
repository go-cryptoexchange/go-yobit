package yobit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	API_BASE    = "https://yobit.net/api"
	API_VERSION = "3"
)

// livecoin represent a livecoin client
type Yobit struct {
	client *client
}

// New returns an instantiated livecoin struct
func New(apiKey, apiSecret string) *Yobit {
	client := NewClient(apiKey, apiSecret)
	return &Yobit{client}
}

// NewWithCustomHttpClient returns an instantiated livecoin struct with custom http client
func NewWithCustomHttpClient(apiKey, apiSecret string, httpClient *http.Client) *Yobit {
	client := NewClientWithCustomHttpConfig(apiKey, apiSecret, httpClient)
	return &Yobit{client}
}

// NewWithCustomTimeout returns an instantiated livecoin struct with custom timeout
func NewWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) *Yobit {
	client := NewClientWithCustomTimeout(apiKey, apiSecret, timeout)
	return &Yobit{client}
}

func (y *Yobit) GetTickers() (map[string]Ticker, error) {
	var tickers TickerInfoResponse
	r, err := y.client.do("GET", "info", "")
	if err != nil {
		return tickers.Pairs, err
	}
	err = json.Unmarshal(r, &tickers)
	if err != nil {
		return tickers.Pairs, err
	}
	return tickers.Pairs, nil
}

func (y *Yobit) GetOrderBook(pair string, limit int) (OrderBook, error) {
	var rawOrderBook map[string]rawOrderBook
	r, err := y.client.do("GET", fmt.Sprintf("depth/%s?limit=%d", pair, limit), "")
	if err != nil {
		return OrderBook{}, err
	}
	err = json.Unmarshal(r, &rawOrderBook)
	if err != nil {
		return OrderBook{}, err
	}

	orderBook := rawOrderBook[pair]
	parsedAsks := make([]Order, 0)
	parsedBids := make([]Order, 0)

	for _, ask := range orderBook.Asks {
		raw := strings.Split(fmt.Sprintf("%v", ask), " ")
		rawQuantity, rawPrice := raw[1], raw[0]
		strQuantity := strings.Trim(rawQuantity, "]")
		strPrice := strings.Trim(rawPrice, "[")
		price, _ := strconv.ParseFloat(strPrice, 64)
		quantity, _ := strconv.ParseFloat(strQuantity, 64)
		parsedAsks = append(parsedAsks, Order{
			Price:    price,
			Quantity: quantity,
		})
	}
	for _, bid := range orderBook.Bids {
		raw := strings.Split(fmt.Sprintf("%v", bid), " ")
		rawQuantity, rawPrice := raw[1], raw[0]
		strQuantity := strings.Trim(rawQuantity, "]")
		strPrice := strings.Trim(rawPrice, "[")
		price, _ := strconv.ParseFloat(strPrice, 64)
		quantity, _ := strconv.ParseFloat(strQuantity, 64)
		parsedBids = append(parsedBids, Order{
			Price:    price,
			Quantity: quantity,
		})
	}

	sort.Slice(parsedAsks, func(i, j int) bool { return parsedAsks[i].Price < parsedAsks[j].Price })
	sort.Slice(parsedBids, func(i, j int) bool { return parsedBids[i].Price > parsedBids[j].Price })

	return OrderBook{
		Asks: parsedAsks,
		Bids: parsedBids,
	}, nil
}
