package coinbase

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

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
	// URL = Coinbase Pro API
	URL = "https://api.pro.coinbase.com/products/BTC-USD/ticker"
)

// GetPrice retrieves Coinbase's price
func GetPrice() float64 {
	res, err := http.Get(URL)
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
