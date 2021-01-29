package main

import (
	"github.com/tsny/btc-alert/binance"
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
		pub := eps.New(ticker.GetPrice, coin, true, 5)
		_ = newListener(pub, conf.Intervals, conf.Thresholds)
		PublisherMap[ticker] = pub
	}

	pub := eps.New(binance.DOGE, "DOGE-USD", true, 30)
	_ = newListener(pub, conf.Intervals, conf.Thresholds)
	PublisherMap["DOGE-USD"] = pub

	utils.Banner("btc-alert initialized")
	for {
	}
}
