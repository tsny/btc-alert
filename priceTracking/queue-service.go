package priceTracking

import "btc-alert/eps"

// todo: probably get rid of these
const (
	SEC_TYPE_STOCK  = "STOCK"
	SEC_TYPE_CRYPTO = "CRYPTO"
)

type Security struct {
	Name   string
	Ticker string
	Type   string
}

type QueueService struct {
	inner map[string]CandleQueue
}

func NewQueueService(publishers ...*eps.Publisher) *QueueService {
	q := QueueService{inner: make(map[string]CandleQueue)}
	return &q
}

func (q *QueueService) TrackSecurities(publishers ...*eps.Publisher) {
	for _, p := range publishers {
		if _, ok := q.inner[p.Ticker]; ok {
			return
		}
		p.Subscribe(q.handlePriceUpdate)
	}
}

func (q *QueueService) handlePriceUpdate(p *eps.Publisher, c eps.Candlestick) {
	queue, ok := q.inner[c.Ticker]
	if !ok {
		queue = NewQueue(10000, c)
		return
	}
	queue.Add(c)
}
