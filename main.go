package main

import (
	"github.com/tsny/btc-alert/binance"
	"github.com/tsny/btc-alert/coinbase"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/utils"
)

var PublisherMap = map[string]*eps.Publisher{}

func main() {
	if conf.DesktopNotifications {
		notif("BTC-ALERT", "Desktop Notifications Enabled", "")
	}

	for _, ticker := range coinbase.CryptoMap {
		pub := eps.New(coinbase.GetPrice, ticker, true, 5)
		_ = newListener(pub, conf.Intervals, conf.Thresholds)
		PublisherMap[ticker] = pub
	}

	pub := eps.New(binance.GetPrice, "DOGEUSDT", true, 30)
	_ = newListener(pub, conf.Intervals, conf.Thresholds)
	PublisherMap["DOGE"] = pub

	utils.Banner("btc-alert initialized")
	for {
	}
}
