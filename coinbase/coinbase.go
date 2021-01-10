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

type result struct {
	TradeID int       `json:"trade_id"`
	Price   string    `json:"price"`
	Size    string    `json:"size"`
	Bid     string    `json:"bid"`
	Ask     string    `json:"ask"`
	Volume  string    `json:"volume"`
	Time    time.Time `json:"time"`
}

const (
	// URL is the Coinbase Pro API base URL
	URL = "https://api.pro.coinbase.com/products/%s/ticker"
	// BTC is Bitcoin
	BTC Source = "BTC-USD"
	// BCH is Bitcoin Cash
	BCH Source = "BCH-USD"
	// Dash is Dash
	Dash Source = "DASH-USD"
	// ETH is Ethereum
	ETH Source = "ETH-USD"
	// EOS is EOS
	EOS Source = "EOS-USD"
	// ETCClassic is Ethererum Classic
	ETCClassic Source = "ETC-USD"
	// ZEC is ZEC
	ZEC Source = "ZEC-USD"
	// MKR is Maker
	MKR Source = "MKR-USD"
	// XLM is Stellar
	XLM Source = "XLM-USD"
	// ATOM is Cosmos
	ATOM Source = "ATOM-USD"
	// LTC is Litecoin
	LTC Source = "LTC-USD"
)

// GetPrice retrieves Coinbase's price
func (s Source) GetPrice() float64 {
	res, err := http.Get(fmt.Sprintf(URL, s))
	if err != nil {
		println(err)
		return -1
	}
	var out result
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		return -1
	}
	p, _ := strconv.ParseFloat(out.Price, 2)
	return p
}
