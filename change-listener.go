package main

import (
	"btc-alert/eps"
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
)

type ChangeListener struct {
	pub *eps.Publisher
}

func NewChangeListener(p *eps.Publisher) *ChangeListener {
	return &ChangeListener{pub: p}
}

func (l *ChangeListener) RegisterPriceMovementListener(userID string, changeAmount float64) {
	startPrice := l.pub.Price()
	id := ""
	logrus.Infof("Alerting %v when %v changes by %v", userID, l.pub.Ticker, changeAmount)
	handler := func(p *eps.Publisher, c *eps.Candlestick, b bool) {
		if diff := math.Abs(c.Price - startPrice); diff > changeAmount {
			msg := fmt.Sprintf("`%v` moved from `%v` to `%v` | A change of `%v`", p.Ticker, startPrice, c.Price, diff)
			_, _ = cryptoBot.SendMessage(msg, userID)
			p.Unsub(id)
		}
	}
	id = l.pub.RegisterPriceUpdateListener(handler)
}

func (l *ChangeListener) RegisterTargetTracker(userID string, target float64) {
	startPrice := l.pub.Price()
	startedBelow := false
	if target > l.pub.Price() {
		startedBelow = true
	}
	id := ""
	logrus.Infof("Alerting %v when %v moves past %v", userID, l.pub.Ticker, target)
	handler := func(p *eps.Publisher, c *eps.Candlestick, b bool) {
		if startedBelow && c.Price > target {
			msg := fmt.Sprintf("`%v` rose above `%v` to `%v`", p.Ticker, startPrice, c.Price)
			_, _ = cryptoBot.SendMessage(msg, userID)
			p.Unsub(id)
		}
		if !startedBelow && c.Price < target {
			msg := fmt.Sprintf("`%v` fell below `%v` to `%v`", p.Ticker, startPrice, c.Price)
			_, _ = cryptoBot.SendMessage(msg, userID)
			p.Unsub(id)
		}
	}
	id = l.pub.RegisterPriceUpdateListener(handler)
}

func (l *ChangeListener) RegisterPercentListener(userID string, percent float64) {
	changeAmount := percent * l.pub.Price()
	l.RegisterPriceMovementListener(userID, changeAmount)
}
