package main

import (
	"btc-alert/eps"
	"fmt"
	"math"

	"github.com/sirupsen/logrus"
)

type ChangeListener struct {
	pub          *eps.Publisher
	startPrice   float64
	changeAmount float64
}

func NewChangeListener(p *eps.Publisher) *ChangeListener {
	return &ChangeListener{pub: p, startPrice: p.Price()}
}

func (l *ChangeListener) Register(userID string, changeAmount float64) {
	l.changeAmount = changeAmount
	id := ""
	logrus.Infof("Alerting %v when %v changes by %v", userID, l.pub.Ticker, changeAmount)
	handler := func(p *eps.Publisher, c *eps.Candlestick, b bool) {
		if diff := math.Abs(c.Price - l.startPrice); diff > l.changeAmount {
			msg := fmt.Sprintf("`%v` moved from `%v` to `%v` | A change of `%v`", p.Ticker, l.startPrice, c.Price, diff)
			_, _ = cryptoBot.SendMessage(msg, userID)
			p.Unsub(id)
		}
	}
	id = l.pub.RegisterPriceUpdateListener(handler)
}

func (l *ChangeListener) RegisterPercentListener(userID string, percent float64) {
	changeAmount := percent * l.pub.Price()
	l.Register(userID, changeAmount)
}
