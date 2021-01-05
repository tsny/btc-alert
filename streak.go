package main

import (
	"encoding/json"
	. "fmt"
	"io/ioutil"
	"math"
	"time"

	"github.com/gen2brain/beeep"
)

type streak struct {
	positive    bool
	occurrences int
	total       float64
}

type interval struct {
	beginPrice       float64
	occurrences      int
	MaxOccurences    int     `json:"maxOccurences"`
	PercentThreshold float64 `json:"percentThreshold"`
	startTime        time.Time
}

var intervals []*interval

func init() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &intervals)
	Printf("Got %d intervals from props file\n", len(intervals))
}

var sf = Sprintf

func onDataUpdated() {
	intervalCompleted := false
	for _, i := range intervals {
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
	s := sf("%s: %s%d Minutes Passed | %s | %s", getTime(), prefix, i.occurrences, totalChange, changes)
	banner(s)
	if math.Abs(percent) > i.PercentThreshold {
		beeep.Alert("BTC_MOVEMENT", s, "assets/warning.png")
	}
}

func getTime() string {
	return time.Now().Format(format)
}

func (i *interval) reset() {
	i.occurrences = 0
	i.startTime = time.Now()
	i.beginPrice = price
}
