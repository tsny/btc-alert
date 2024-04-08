package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
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
	BTC  = "BTC-USD"
	ETH  = "ETH-USD"
	DOGE = "DOGE-USD"
	SOL  = "SOL-USD"
)

func Tickers() []string {
	return []string{
		BTC, ETH, DOGE, SOL,
	}
}

func FindTicker(ticker string) (string, bool) {
	ticker = strings.ToLower(ticker)
	return lo.Find(Tickers(), func(e string) bool {
		return strings.Contains(strings.ToLower(e), ticker)
	})
}

// TickerURL is the Coinbase Ticker API URL
// https://api.pro.coinbase.com/products/BTC-USD/ticker
const TickerURL = "https://api.pro.coinbase.com/products/%s/ticker"

// DailyURL is the Coinbase Ticker API URL that returns the stats for the last 24h
const DailyURL = "https://api.pro.coinbase.com/products/%s/stats"

// GetDetails retrieves details regarding a specific coin
func GetDetails(ticker string) (*Ticker, error) {
	res, err := http.Get(fmt.Sprintf(TickerURL, ticker))
	if err != nil {
		return nil, err
	}
	var out Ticker
	if err = json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetPrice retrieves Coinbase's price for a specific coin
func GetPrice(ticker string) float64 {
	out, err := GetDetails(ticker)
	if err != nil {
		println(err.Error())
		return -1
	}
	p, _ := strconv.ParseFloat(out.Price, 64)
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
