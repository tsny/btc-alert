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
	url    = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=BTC-USD"
	up     = "🟩"
	down   = "🟥"
	format = "15:04:05"
)

var lastPrice = 0.00
var price = 0.00
var priceAtLastHourMark = 0.00
var priceAtLast5MinMark = 0.00
var timeOfLastHourMark time.Time
var timeOfLast5MinMark time.Time
var minutesPassed = 0
var streak Streak

func main() {
	banner("Fetching BTC Prices...")
	for {
		price = fetchData()
		if minutesPassed == 0 {
			priceAtLastHourMark = price
			priceAtLast5MinMark = price
			timeOfLastHourMark = time.Now()
			timeOfLast5MinMark = time.Now()
		}
		diff := price - lastPrice
		emoji := getEmoji()
		t := time.Now().Format(format)

		if lastPrice == 0.00 {
			fmt.Printf("%s %s: $%.2f\n", emoji, t, price)
		} else {
			fmt.Printf("%s %s: $%v | Change: $%.2f\n", emoji, t, price, diff)
		}

		lastPrice = price
		time.Sleep(60 * time.Second)

		minutesPassed++
		if minutesPassed > 59 {
			onHourPassed()
		} else if minutesPassed%5 == 0 {
			onFiveMinPassed()
		}
	}
}

func banner(str string) {
	b := strings.Repeat("-", len(str))
	fmt.Printf("%s\n%s\n%s\n", b, str, b)
}

func onFiveMinPassed() {
	timeOfLast5MinMark = time.Now()
	banner(fmt.Sprintf("5 Minutes Passed - Change: $%.2f", price-priceAtLast5MinMark))
	priceAtLast5MinMark = price
}

func onHourPassed() {
	timeOfLastHourMark = time.Now()
	minutesPassed = 0
	banner(fmt.Sprintf("Hour Passed - Change: $%.2f", price-priceAtLastHourMark))
	priceAtLastHourMark = price
}

func getEmoji() string {
	if price > lastPrice {
		return up
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
