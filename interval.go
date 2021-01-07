package main

import (
	"fmt"

	"github.com/tsny/btc-alert/eps"
)

type listener struct {
	intervals  []*interval
	thresholds []*threshold
}

func newListener(p *eps.Publisher, intervals []*interval, thresholds []*threshold) *listener {
	i := &listener{intervals, thresholds}
	p.Subscribe(i.onPriceUpdated)
	return i
}

func (i *listener) onPriceUpdated(new, old float64) {
	i.checkIntervals(new, old)
	i.checkThresholds(new, old)
	fmt.Print(getSummaryNew(new, old))
	if conf.Discord.Enabled {
		discordMessage(getSummaryNew(new, old), false)
	}
}

func (i *listener) checkIntervals(new, old float64) {
	for _, i := range i.intervals {
		if i.beginPrice == 0 {
			i.beginPrice = new
		}
		i.occurrences++
		if i.occurrences >= i.MaxOccurences {
			i.onCompleted(new, old)
			i.reset(new)
		}
	}
}

func (i *listener) checkThresholds(new, old float64) {
	for _, t := range i.thresholds {
		if t.beginPrice == 0 {
			t.beginPrice = new
		}
		if new >= t.Threshold+t.beginPrice {
			t.onThresholdReached(true, new, old)
		} else if new <= t.beginPrice-t.Threshold {
			t.onThresholdReached(false, new, old)
		}
	}
}

func getSummaryNew(new, old float64) string {
	t := getTime()
	emoji := getEmoji(new, old)
	diff := new - old
	percent := (diff / new) * 100
	if old == 0.00 {
		return sf("%s %s: $%.2f \n", emoji, t, new)
	}
	return sf("%s %s: $%.2f | Chg: $%.2f | Percent: %.3f%% \n", emoji, t, new, diff, percent)
}
