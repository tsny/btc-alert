package main

import (
	"encoding/json"
	. "fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/gen2brain/beeep"
)

type config struct {
	Intervals []*interval `json:"intervals"`
	Changes   []*change   `json:"changes"`
}

// changes are price jumps
// alerts after prices move a certain amount from
// the starting price
type change struct {
	beginPrice float64
	Threshold  float64 `json:"threshold"`
}

// intervals
type interval struct {
	beginPrice       float64
	occurrences      int
	MaxOccurences    int     `json:"maxOccurences"`
	PercentThreshold float64 `json:"percentThreshold"`
	startTime        time.Time
}

var conf config

func init() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &conf)
	Printf("props: %d intervals | %d changes\n", len(conf.Intervals), len(conf.Changes))
}

var sf = Sprintf

func onDataUpdated() {
	intervalCompleted := false
	for _, i := range conf.Intervals {
		if i.beginPrice == 0 {
			i.beginPrice = price
		}
		i.occurrences++
		if i.occurrences >= i.MaxOccurences {
			i.onCompleted()
			i.reset()
			intervalCompleted = true
		}
	}
	if !intervalCompleted {
		t := getTime()
		emoji := getEmoji(price, lastPrice)
		diff := price - lastPrice
		percent := (diff / price) * 100
		if lastPrice == 0.00 {
			Printf("%s %s: $%.2f \n", emoji, t, price)
		} else {
			Printf("%s %s: $%.2f | Change: $%.2f | Percent: %.3f%% \n", emoji, t, price, diff, percent)
		}
	}
}

func (i *interval) onCompleted() {
	diff := price - i.beginPrice
	percent := (diff / i.beginPrice) * 100
	prefix := ""
	if math.Abs(percent) > i.PercentThreshold {
		prefix = sf("%s ALERT: Threshold of %.2f%% was reached! ", alert, i.PercentThreshold)
	}
	totalChange := sf("$%.2f --> $%.2f", i.beginPrice, price)
	changes := sf("Change: $%.2f | Percent: %.3f%%", diff, percent)
	bannerText := sf("%s: %s%d Minutes Passed | %s | %s", getTime(), prefix, i.occurrences, totalChange, changes)
	banner(bannerText)
	if math.Abs(percent) > i.PercentThreshold {
		hdr := sf("%d Minutes Passed | %.2f%%", i.MaxOccurences, i.PercentThreshold)
		beeep.Alert(hdr, bannerText, "assets/warning.png")
	}
}

func (c *change) checkThreshold() {
	if price >= c.Threshold+c.beginPrice {
		c.onThresholdReached()
	} else if price <= c.beginPrice-c.Threshold {
		c.onThresholdReached()
	}
}

func (c *change) onThresholdReached() {
	c.beginPrice = price
	str := "%s: ALERT: Price Threshold Breached: %d | %s"
	bannerf(str, getTime(), c.Threshold, formatPriceMovement(c.beginPrice, price))
}

func formatPriceMovement(begin, end float64) string {
	return sf("$%.2f --> $%.2f", begin, end)
}

func getTime() string {
	return time.Now().Format(format)
}

func (i *interval) reset() {
	i.occurrences = 0
	i.startTime = time.Now()
	i.beginPrice = price
}
