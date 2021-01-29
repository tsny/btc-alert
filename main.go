package main

import (
	"github.com/tsny/btc-alert/coinbase"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/utils"
	"github.com/tsny/btc-alert/yahoo"
)

var PublisherMap = map[coinbase.Source]*eps.Publisher{}

func main() {
	if conf.DesktopNotifications {
		notif("BTC-ALERT", "Desktop Notifications Enabled", "")
	}

	for coin, ticker := range coinbase.CryptoMap {
		pub := eps.New(ticker.GetPrice, coin, true, 5)
		_ = newListener(pub, conf.Intervals, conf.Thresholds)
		PublisherMap[ticker] = pub
	}

	s := yahoo.Source("DOGE-USD")
	pub := eps.New(s.GetPrice, "DOGE-USD", true, 30)
	_ = newListener(pub, conf.Intervals, conf.Thresholds)
	PublisherMap["DOGE-USD"] = pub

	utils.Banner("btc-alert initialized")
	for {
	}
}
