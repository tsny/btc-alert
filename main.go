package main

import (
	"github.com/tsny/btc-alert/coinbase"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/utils"
)

var PublisherMap = map[coinbase.Source]*eps.Publisher{}

func main() {
	if conf.DesktopNotifications {
		notif("BTC-ALERT", "Desktop Notifications Enabled", "")
	}

	for coin, ticker := range coinbase.CryptoMap {
		pub := eps.New(ticker.GetPrice, coin, true)
		_ = newListener(pub, conf.Intervals, conf.Thresholds)
		PublisherMap[ticker] = pub
	}

	utils.Banner("btc-alert initialized")
	for {
	}
}
