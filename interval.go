package main

import (
	"fmt"

	"github.com/tsny/btc-alert/eps"
)

type listener struct {
	intervals  []interval
	thresholds []threshold
}

func newListener(p *eps.Publisher, intervals []interval, thresholds []threshold) *listener {
	i := &listener{intervals, thresholds}
	p.Subscribe(i.onPriceUpdated)
	return i
}

func (i *listener) onPriceUpdated(p *eps.Publisher, new, old float64) {
	i.checkIntervals(p, new, old)
	i.checkThresholds(p, new, old)

	fmt.Print(getSummaryNew(p, new, old))
	// if conf.Discord.Enabled {
	// 	discordMessage(getSummaryNew(p, new, old), false)
	// }
}

func (i *listener) checkIntervals(p *eps.Publisher, new, old float64) {
	for _, i := range i.intervals {
		if i.beginPrice == 0 {
			i.beginPrice = new
		}
		i.occurrences++
		if i.occurrences >= i.MaxOccurences {
			i.onCompleted(p, new, old)
			i.reset(new)
		}
	}
}

func (i *listener) checkThresholds(p *eps.Publisher, new, old float64) {
	for _, t := range i.thresholds {
		if t.beginPrice == 0 {
			t.beginPrice = new
		}
		if new >= t.Threshold+t.beginPrice {
			t.onThresholdReached(p, true, new, old)
		} else if new <= t.beginPrice-t.Threshold {
			t.onThresholdReached(p, false, new, old)
		}
	}
}

func getSummaryNew(p *eps.Publisher, new, old float64) string {
	t := getTime()
	emoji := getEmoji(new, old)
	diff := new - old
	percent := (diff / new) * 100
	if old == 0.00 {
		return sf("%s %s: (%s) $%.2f \n", emoji, t, p.Source, new)
	}
	s := "%s %s: (%s) $%.2f | Chg: $%.2f | Percent: %.3f%% \n"
	return sf(s, emoji, t, p.Source, new, diff, percent)
}
