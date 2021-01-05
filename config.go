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
}

// changes are price jumps
// alerts after prices move a certain amount from
// the starting price
type threshold struct {
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
	fmt.Printf("props: %d intervals | %d thresholds\n", len(conf.Intervals), len(conf.Thresholds))
}
