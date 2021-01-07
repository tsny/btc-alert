package coindesk

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// URL = Coindesk API - Current Price endpoint
const URL = "https://api.coindesk.com/v1/bpi/currentprice.json"

type Result struct {
	Time       Time   `json:"time"`
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpi        Bpi    `json:"bpi"`
}

type Time struct {
	Updated    string    `json:"updated"`
	UpdatedISO time.Time `json:"updatedISO"`
	Updateduk  string    `json:"updateduk"`
}

type USD struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

type GBP struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

type EUR struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

type Bpi struct {
	USD USD `json:"USD"`
	GBP GBP `json:"GBP"`
	EUR EUR `json:"EUR"`
}

// GetPrice retrieves Coindesk's price
func GetPrice() float64 {
	res, err := http.Get(URL)
	if err != nil {
		println(err)
		return -1
	}
	var out Result
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		panic(err)
	}
	s := strings.ReplaceAll(out.Bpi.USD.Rate, ",", "")
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		println(err.Error())
		return -1
	}
	return price
}
