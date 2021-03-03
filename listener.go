package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tsny/btc-alert/eps"
)

type listener struct {
	publisher  *eps.Publisher
	lastAlert  time.Time
	intervals  []interval
	thresholds []threshold
}

func newListener(p *eps.Publisher, intervals []interval, thresholds []threshold) *listener {
	// TODO: rethink this and get rid of this garbage
	cpi := intervals
	cpt := thresholds
	l := listener{}
	for _, i := range cpi {
		l.intervals = append(l.intervals, i)
	}
	for _, t := range cpt {
		l.thresholds = append(l.thresholds, t)
	}
	p.Subscribe(l.onPriceUpdated)
	l.publisher = p
	return &l
}

func (l *listener) onPriceUpdated(p *eps.Publisher, c eps.Candlestick) {
	if !c.Complete || c.Current == 0 {
		return
	}
	if p.Streak > conf.StreakAlert && conf.StreakAlert > 0 {
		if conf.Discord.Enabled {
			str := fmt.Sprintf("%s has a streak of %v", p.Ticker, p.Streak)
			cryptoBot.SendMessage(str, "", false)
		}
	}
	go l.checkIntervals(p, c.Current, c.Previous)
	go l.checkThresholds(p, c.Current, c.Previous)

	s := c.String()
	if c.Volatility() > conf.VolatilityAlert {
		cryptoBot.SendMessage(s, "everyone", false)
		s = s + " <-- ALERT"
	}
	log.Print(s)
}

func (l *listener) checkIntervals(p *eps.Publisher, new, old float64) {
	for i, interval := range l.intervals {
		if interval.beginPrice == 0 {
			l.intervals[i].beginPrice = new
		}
		// fmt.Printf("%d minute interval %s occurred %d times\n",
		// 	interval.MaxOccurences, p.Source, interval.occurrences)
		l.intervals[i].occurrences++
		if interval.occurrences >= interval.MaxOccurences {
			interval.onCompleted(p, new, old)
			l.intervals[i].reset(new)
		}
	}
}

func (l *listener) checkThresholds(p *eps.Publisher, new, old float64) {
	for i, t := range l.thresholds {
		if t.beginPrice == 0 {
			l.thresholds[i].beginPrice = new
			continue
		}
		if new >= t.Threshold+t.beginPrice {
			l.thresholds[i].onThresholdReached(p, true, new, old)
		} else if new <= t.beginPrice-t.Threshold {
			l.thresholds[i].onThresholdReached(p, false, new, old)
		}
	}
}
