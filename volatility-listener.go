package main

import (
	"btc-alert/eps"
	"btc-alert/utils"

	"github.com/sirupsen/logrus"
)

type VolatilityListener struct {
	pub           *eps.Publisher
	startPrice    float64
	percentChange float64
	numCandles    int
}

func NewVolatilityListener(p *eps.Publisher, percentChange float64, durInMinutes int) *VolatilityListener {
	if percentChange > 1 {
		percentChange *= .01
	}
	v := &VolatilityListener{
		pub:           p,
		startPrice:    p.Price(),
		percentChange: percentChange,
		numCandles:    durInMinutes,
	}
	p.RegisterPriceUpdateListener(v.HandlePriceUpdate)
	return v
}

func (v *VolatilityListener) HandlePriceUpdate(p *eps.Publisher, c *eps.Candlestick, complete bool) {
	if !complete {
		return
	}
	if len(p.Stack.Array) < v.numCandles {
		return
	}
	arr := p.Stack.Array
	last := arr[v.numCandles-1]
	logrus.Infof("%v | %v => %v (%.2f)", utils.CompareTimes(c.Start, last.Start), last.Price, arr[0].Price, c.DiffPercent(last))
}
