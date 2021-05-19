package priceTracking

import (
	"btc-alert/eps"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Todo: this service should also use security type instead of string
type QueueService struct {
	inner map[string]*CandleQueue
	mutex *sync.Mutex
}

// NewQueueService is a constructor
func NewQueueService(publishers ...*eps.Publisher) *QueueService {
	q := QueueService{make(map[string]*CandleQueue), &sync.Mutex{}}
	q.TrackSecurities(publishers...)
	return &q
}

// Creates candle queues for every publisher passed in
func (q *QueueService) TrackSecurities(publishers ...*eps.Publisher) {
	for _, p := range publishers {
		if _, ok := q.inner[p.Ticker]; ok {
			return
		}
		p.RegisterSubscriber(q.handlePriceUpdate)
	}
}

// Finds a candle queue by its ticker
func (q *QueueService) FindByTicker(ticker string) *CandleQueue {
	if queue, ok := q.inner[ticker]; ok {
		return queue
	}
	return nil
}

func (q *QueueService) handlePriceUpdate(p *eps.Publisher, c eps.Candlestick) {
	// Ignore incomplete candles
	if !c.Complete {
		return
	}
	queue, ok := q.inner[c.Ticker]
	if !ok {
		q.mutex.Lock()
		q.inner[c.Ticker] = NewQueue(capLimit, c)
		q.mutex.Unlock()
		log.Infof("Creating new queue for %s", p.String())
		return
	}
	queue.Add(c)
}

// Returns all queues within the service as an array
func (q *QueueService) AllQueues() []*CandleQueue {
	var arr []*CandleQueue
	for _, v := range q.inner {
		arr = append(arr, v)
	}
	return arr
}
