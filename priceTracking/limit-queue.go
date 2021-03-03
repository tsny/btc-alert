package priceTracking

import "btc-alert/eps"

const capLimit = 10000

type CandleQueue struct {
	Oldest eps.Candlestick
	Newest eps.Candlestick
	inner  []eps.Candlestick
	cap    int
}

func NewQueue(cap int, data ...eps.Candlestick) CandleQueue {
	q := CandleQueue{cap: cap}
	if len(data) > 0 {
		q.inner = data
		q.Newest = data[0]
	}
	return q
}

func (l *CandleQueue) Add(c eps.Candlestick) {
	// todo: make this better than just truncating
	if len(l.inner) > l.cap {
		l.inner = l.inner[:capLimit-1]
	}
	l.inner = append(l.inner, c)
	l.Newest = c
}

func (l *CandleQueue) GetCap() int {
	return l.cap
}

func (l *CandleQueue) SetCap(cap int) {
	if cap > capLimit {
		cap = capLimit
	}
	l.cap = cap
}

func (l *CandleQueue) GetQueue() []eps.Candlestick {
	return l.inner
}
