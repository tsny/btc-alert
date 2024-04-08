package binance

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const URL = "https://api.binance.com/api/v3/ticker/price?symbol=%s"

type result struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

// GetPrice gets the price of the asset in Binance
func GetPrice(ticker string) float64 {
	res, err := http.Get(fmt.Sprintf(URL, ticker))
	if err != nil {
		println(err)
		return -1
	}
	var out result
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	price, err := strconv.ParseFloat(out.Price, 64)
	if err != nil {
		println(err.Error())
		return -1
	}
	return price
}

func DOGE() float64 {
	return GetPrice("DOGEUSDT")
}
