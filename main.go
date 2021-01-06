package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

const (
	// url     = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=BTC-USD"
	url     = "https://api.coindesk.com/v1/bpi/currentprice.json"
	up      = "🟩"
	down    = "🟥"
	neutral = "🟦"
	alert   = "☎️"
	dollar  = "💲"
	format  = "03:04 PM"
)

var lastPrice = 0.00
var price = 0.00

func main() {
	time.Sleep(2 * time.Second)
	if conf.BootNotification {
		beeep.Alert("BTC_ALERT", "STARTING UP", "assets/warning.png")
	}
	banner("btc-alert initialized")
	for {
		price = fetchData()
		onPriceUpdated()
		lastPrice = price
		time.Sleep(60 * time.Second)
	}
}

func onFirstPriceFetched() {
	for _, c := range conf.Thresholds {
		c.beginPrice = price
	}
}

func getEmoji(curr, prev float64) string {
	if prev < curr {
		return up
	} else if prev == curr {
		return neutral
	}
	return down
}

func fetchData() float64 {
	res, err := http.Get(url)
	if err != nil {
		banner(err.Error())
		return 0
	}
	var out CoindeskResponse
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		panic(err)
	}
	// price = out.QuoteResponse.Result[0].RegularMarketPrice
	s := strings.ReplaceAll(out.Bpi.USD.Rate, ",", "")
	price, err = strconv.ParseFloat(s, 64)
	if err != nil {
		println(err.Error())
		return -1
	}
	if lastPrice == 0 {
		onFirstPriceFetched()
	}
	return math.Round(price*100) / 100
}
