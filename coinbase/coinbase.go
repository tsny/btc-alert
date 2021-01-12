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

// URL is the Coinbase Ticker API URL
const URL = "https://api.pro.coinbase.com/products/%s/ticker"

// GetPrice retrieves Coinbase's price
func (s Source) GetPrice() float64 {
	res, err := http.Get(fmt.Sprintf(URL, s))
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
