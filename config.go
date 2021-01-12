package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tsny/btc-alert/utils"
)

type config struct {
	Intervals            []interval    `json:"intervals"`
	Thresholds           []threshold   `json:"thresholds"`
	VolatilityAlert      float64       `json:"volatilityAlert"`
	BootNotification     bool          `json:"bootNotification"`
	DesktopNotifications bool          `json:"desktopNotifications"`
	Discord              discordConfig `json:"discord"`
}

type discordConfig struct {
	Token     string `json:"token"`
	Enabled   bool   `json:"enabled"`
	ChannelID string `json:"channelId"`
}

var conf config

func init() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &conf)
	utils.Banner("btc-alert initializing")
	fmt.Printf("props: %d intervals | %d thresholds\n", len(conf.Intervals), len(conf.Thresholds))
	for _, i := range conf.Intervals {
		fmt.Printf("Interval -- Minutes: %d | Percentage Threshold: %v\n", i.MaxOccurences, i.PercentThreshold)
	}
	if conf.Discord.Enabled {
		println("Discord enabled")
		cryptoBot = initBot(conf.Discord.Token)
	}
}
