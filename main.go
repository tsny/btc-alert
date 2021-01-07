package main

import (
	"github.com/tsny/btc-alert/eps"
)

func main() {
	publisher := eps.New(eps.CoindeskURL)
	_ = newListener(publisher, conf.Intervals, conf.Thresholds)
	banner("btc-alert initialized")
	publisher.StartListening()
}
