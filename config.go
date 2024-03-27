package main

import (
	"encoding/json"
	"os"

	"btc-alert/utils"
)

type config struct {
	PercentageChanges []PercentChange `json:"percentageChanges"`

	VolatilityAlert      float64       `json:"volatilityAlert"`
	BootNotification     bool          `json:"bootNotification"`
	DesktopNotifications bool          `json:"desktopNotifications"`
	Discord              discordConfig `json:"discord"`
	StreakAlert          int           `json:"streakAlert"`
}

type discordConfig struct {
	Token                        string   `json:"token"`
	Enabled                      bool     `json:"enabled"`
	ChannelID                    string   `json:"channelId"`
	MessageForEachIntervalUpdate bool     `json:"messageOnIntervals"`
	AlertEveryone                bool     `json:"alertEveryone"` // Whether to tag @everyone when an alert is sent
	UsersToNotify                []string `json:"usersToNotify"`
}

var conf config

func readConfig() {
	cfgPath := utils.Getenv("BTC_ALERT_CONFIG_PATH", "config.json")
	println(cfgPath)
	bytes, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(bytes, &conf); err != nil {
		panic(err)
	}
	utils.Banner("btc-alert initializing")

	if conf.Discord.Enabled {
		cryptoBot = initBot(conf.Discord.Token)
	}
}
