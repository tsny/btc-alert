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
	var cpi []interval
	var cpt []threshold
	copy(intervals, cpi)
	copy(thresholds, cpt)
	l := listener{cpi, cpt}
	p.Subscribe(l.onPriceUpdated)
	return &l
}

func (l *listener) onPriceUpdated(p *eps.Publisher, new, old float64) {
	l.checkIntervals(p, new, old)
	l.checkThresholds(p, new, old)

	fmt.Print(getSummaryNew(p, new, old))
	// if conf.Discord.Enabled {
	// 	discordMessage(getSummaryNew(p, new, old), false)
	// }
}

func (l *listener) checkIntervals(p *eps.Publisher, new, old float64) {
	for i, interval := range l.intervals {
		if interval.beginPrice == 0 {
			l.intervals[i].beginPrice = new
		}
		l.intervals[i].occurrences++
		if interval.occurrences >= interval.MaxOccurences {
			interval.onCompleted(p, new, old)
			interval.reset(new)
		}
	}
}

func (l *listener) checkThresholds(p *eps.Publisher, new, old float64) {
	for i, t := range l.thresholds {
		if t.beginPrice == 0 {
			l.thresholds[i].beginPrice = new
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
