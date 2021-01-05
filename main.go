package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	url     = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=BTC-USD"
	up      = "🟩"
	down    = "🟥"
	neutral = "🟦"
	alert   = "☎️"
	format  = "15:04:05"
)

var lastPrice = 0.00
var price = 0.00

func main() {
	banner("Fetching BTC Prices...")
	for {
		price = fetchData()
		onDataUpdated()
		lastPrice = price
		time.Sleep(60 * time.Second)
	}
}

func banner(str string, args ...interface{}) {
	str = sf(str, args...)
	b := strings.Repeat("-", len(str))
	fmt.Printf("%s\n%s\n%s\n", b, str, b)
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
		panic(err)
	}
	var out TLR
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		panic(err)
	}
	price := out.QuoteResponse.Result[0].RegularMarketPrice
	return math.Round(price*100) / 100
}
