package alert

import (
	"btc-alert/pkg/eps"
	"fmt"
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
	candles := p.Stack.Array
	if len(candles) < v.numCandles {
		return
	}
	startCandle := candles[v.numCandles-1]
	for i := v.numCandles - 2; i >= 0; i-- {
		currCandle := candles[i]
		diffPercent := currCandle.DiffPercent(startCandle)
		if diffPercent > v.percentChange {
			// TODO: DO THE NOTIF
		}
		fmt.Println(currCandle.DiffString(*startCandle))
	}
}
