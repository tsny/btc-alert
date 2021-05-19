package coinbase

import (
	"btc-alert/eps"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

// CryptoMap is a map of the Coin's simple name to its ticker
var CryptoMap = []*eps.Security{
	{Name: "Bitcoin", Ticker: "BTC-USD"},
	{Name: "Bitcoin Cash", Ticker: "BCH-USD"},
	{Name: "DASH", Ticker: "DASH-USD"},
	{Name: "Ethereum", Ticker: "ETH-USD"},
	{Name: "EOS", Ticker: "EOS-USD"},
	{Name: "Ethereum Classic", Ticker: "ETC-USD"},
	{Name: "ZEC", Ticker: "ZEC-USD"},
	{Name: "Maker", Ticker: "MKR-USD"},
	{Name: "XLM", Ticker: "XLM-USD"},
	{Name: "ATOM", Ticker: "ATOM-USD"},
	{Name: "Lite Coin", Ticker: "LTC-USD"},
}

func init() {
	for _, c := range CryptoMap {
		c.AdditionalNames = append(c.AdditionalNames, strings.ReplaceAll(c.Ticker, "-USD", ""))
		c.Type = eps.Crypto
		c.Source = "Coinbase"
	}
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
	res, err := http.Get(fmt.Sprintf(DailyURL, ticker))
	if err != nil {
		println(err)
		return nil
	}
	var out OneDayCandle
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		return nil
	}
	return &out
}
