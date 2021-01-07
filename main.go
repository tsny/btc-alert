package main

import (
	"github.com/tsny/btc-alert/eps"
)

func main() {
	if conf.DesktopNotifications {
		notif("BTC-ALERT", "Desktop Notifications Enabled", "")
	}
	publisher := eps.New(eps.CoindeskURL)
	_ = newListener(publisher, conf.Intervals, conf.Thresholds)
	banner("btc-alert initialized")
	publisher.StartListening()
}
