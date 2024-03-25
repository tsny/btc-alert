package coinbase

import (
	"btc-alert/eps"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Ticker struct {
	TradeID int       `json:"trade_id"`
	Price   string    `json:"price"`
	Size    string    `json:"size"`
	Bid     string    `json:"bid"`
	Ask     string    `json:"ask"`
	Volume  string    `json:"volume"`
	Time    time.Time `json:"time"`
}

type OneDayCandle struct {
	Open        string `json:"open"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Volume      string `json:"volume"`
	Last        string `json:"last"`
	Volume30Day string `json:"volume_30day"`
}

const (
	BTC = "BTC-USD"
	ETH = "ETH-USD"
)

// CryptoMap is a map of the Coin's simple name to its ticker
var CryptoMap = map[string]*eps.Security{
	BTC: {Name: "BTC", Type: eps.Crypto},
}

// TickerURL is the Coinbase Ticker API URL
// https://api.pro.coinbase.com/products/BTC-USD/ticker
const TickerURL = "https://api.pro.coinbase.com/products/%s/ticker"

// DailyURL is the Coinbase Ticker API URL that returns the stats for the last 24h
const DailyURL = "https://api.pro.coinbase.com/products/%s/stats"

// GetDetails retrieves details regarding a specific coin
func GetDetails(ticker string) *Ticker {
	res, err := http.Get(fmt.Sprintf(TickerURL, ticker))
	if err != nil {
		println(err)
		return nil
	}
	var out Ticker
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		return nil
	}
	return &out
}

// GetPrice retrieves Coinbase's price for a specific coin
func GetPrice(ticker string) float64 {
	out := GetDetails(ticker)
	if out == nil {
		return -1
	}
	p, _ := strconv.ParseFloat(out.Price, 2)
	return p
}

// Get24Hour returns a 24 hour candlestick for a ticker
func Get24Hour(ticker string) *OneDayCandle {
	url := fmt.Sprintf(DailyURL, ticker)
	println(url)
	res, err := http.Get(url)
	if err != nil {
		println(err)
		return nil
	}
	var out OneDayCandle
	d := json.NewDecoder(res.Body)
	if err = d.Decode(&out); err != nil {
		return nil
	}
	return &out
}
