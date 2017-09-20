package yobit

import (
	"testing"
)

func TestGetTickers(t *testing.T) {
	api := New("", "")
	tickers, err := api.GetTickers()
	if err != nil {
		t.Fatal(err)
	}
	if len(tickers) < 1 {
		t.Fail()
	}
	if ticker, ok := tickers["btc_usd"]; !ok {
		t.Fatal(ticker)
	}
	for k, ticker := range tickers {
		if ticker.Hidden != 0 {
			t.Log(k, ticker)
		}
	}
}

func TestGetOrderBook(t *testing.T) {
	api := New("", "")
	orderBook, err := api.GetOrderBook("btc_usd", 150)
	if err != nil {
		t.Fatal(err)
	}
	if len(orderBook.Asks) < 1 {
		t.Fail()
	}
}
