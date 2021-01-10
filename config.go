package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type config struct {
	Intervals            []interval    `json:"intervals"`
	Thresholds           []threshold   `json:"thresholds"`
	BootNotification     bool          `json:"bootNotification"`
	DesktopNotifications bool          `json:"desktopNotifications"`
	Discord              discordConfig `json:"discord"`
}

type discordConfig struct {
	Token              string `json:"token"`
	Enabled            bool   `json:"enabled"`
	ChannelID          string `json:"channelId"`
	ClearChannelOnBoot bool   `json:"clearOnBoot"`
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
	if conf.Discord.Enabled {
		println("Discord enabled")
		go initBot(conf.Discord.Token)
	}
}
