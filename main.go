package main

import (
	"github.com/tsny/btc-alert/coinbase"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/utils"
)

func main() {
	if conf.DesktopNotifications {
		notif("BTC-ALERT", "Desktop Notifications Enabled", "")
	}

	btc := eps.New(coinbase.BTC.GetPrice, "BTC")
	_ = newListener(btc, conf.Intervals, conf.Thresholds)

	dash := eps.New(coinbase.Dash.GetPrice, "DASH")
	_ = newListener(dash, conf.Intervals, conf.Thresholds)

	bch := eps.New(coinbase.BCH.GetPrice, "BCH")
	_ = newListener(bch, conf.Intervals, conf.Thresholds)

	ltc := eps.New(coinbase.LTC.GetPrice, "LTC")
	_ = newListener(ltc, conf.Intervals, conf.Thresholds)

	eth := eps.New(coinbase.ETH.GetPrice, "ETH")
	_ = newListener(eth, conf.Intervals, conf.Thresholds)

	etc := eps.New(coinbase.ETCClassic.GetPrice, "ETC")
	_ = newListener(etc, conf.Intervals, conf.Thresholds)

	eos := eps.New(coinbase.EOS.GetPrice, "EOS")
	_ = newListener(eos, conf.Intervals, conf.Thresholds)

	dash.StartProducing()
	btc.StartProducing()
	bch.StartProducing()
	ltc.StartProducing()
	eth.StartProducing()
	etc.StartProducing()
	eos.StartProducing()

	utils.Banner("btc-alert initialized")
	for {
	}
}
