package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Source is a type alias for strings
type Source string

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
var CryptoMap = map[string]Source{
	"BTC":  "BTC-USD",
	"BCH":  "BCH-USD",
	"DASH": "DASH-USD",
	"ETH":  "ETH-USD",
	"EOS":  "EOS-USD",
	"ETC":  "ETC-USD",
	"ZEC":  "ZEC-USD",
	"MKR":  "MKR-USD",
	"XLM":  "XLM-USD",
	"ATOM": "ATOM-USD",
	"LTC":  "LTC-USD",
}

// TickerURL is the Coinbase Ticker API URL
const TickerURL = "https://api.pro.coinbase.com/products/%s/ticker"

// DailyURL is the Coinbase Ticker API URL that returns the stats for the last 24h
const DailyURL = "https://api.pro.coinbase.com/products/%s/stats"

// GetPrice retrieves Coinbase's price
func (s Source) GetPrice() float64 {
	res, err := http.Get(fmt.Sprintf(TickerURL, s))
	if err != nil {
		println(err)
		return -1
	}
	var out Ticker
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		return -1
	}
	p, _ := strconv.ParseFloat(out.Price, 2)
	return p
}

func (s Source) Get24Hour() *OneDayCandle {
	res, err := http.Get(fmt.Sprintf(DailyURL, s))
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
