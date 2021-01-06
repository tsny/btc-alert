package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type config struct {
	Intervals        []*interval  `json:"intervals"`
	Thresholds       []*threshold `json:"thresholds"`
	BootNotification bool         `json:"bootNotification"`
	UseDiscord       bool         `json:"useDiscord"`
}

// thresholds are price jumps
// they alert after prices move a certain amount from
// the starting price
type threshold struct {
	beginPrice float64
	Threshold  float64 `json:"threshold"`
}

// intervals are checked every minute
// if 'maxOccurences' number of minutes pass
// then the interval lapses and onCompleted() is called
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
	banner("btc-alert initializing")
	fmt.Printf("props: %d intervals | %d thresholds\n", len(conf.Intervals), len(conf.Thresholds))
	for _, i := range conf.Intervals {
		fmt.Printf("Interval -- Minutes: %d | Percentage Threshold: %v\n", i.MaxOccurences, i.PercentThreshold)
	}
	if conf.UseDiscord {
		token := initToken()
		go initBot(token)
	}
}
