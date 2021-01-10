// Package eps stands for Exchange Publisher Service
package eps

import (
	"fmt"
	"time"
)

const (
	// YahooURL = Yahoo Finance
	YahooURL = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=BTC-USD"
)

// Publisher periodically grabs data from its URL
// and sends out updates with the price it gets back
type Publisher struct {
	Source          string
	priceUpdateSubs []func(*Publisher, float64, float64)
	price           float64
	lastPrice       float64
	active          bool
	sleepDuration   int
	priceFetcher    func() float64
}

// New is a constructor
func New(priceFetcher func() float64, source string) *Publisher {
	return &Publisher{
		source,
		[]func(p *Publisher, new, old float64){},
		0,
		0,
		false,
		60,
		priceFetcher,
	}
}

// StartListening loops and updates the price from the chosen exchange
func (p *Publisher) StartListening() {
	fmt.Printf("%s -- Price Event Publisher active\n", p.Source)
	if p.active {
		return
	}
	p.active = true
	go func() {
		for {
			p.fetchAndUpdatePrice()
			time.Sleep(time.Duration(p.sleepDuration) * time.Second)
		}
	}()
}

// Subscribe allows services to subscribe to new BitCoin events
func (p *Publisher) Subscribe(f func(p *Publisher, new, old float64)) {
	fmt.Printf("%s Publisher has new subscriber\n", p.Source)
	p.priceUpdateSubs = append(p.priceUpdateSubs, f)
}

func (p *Publisher) onPriceUpdated() {
	for _, c := range p.priceUpdateSubs {
		c(p, p.price, p.lastPrice)
	}
}

func (p *Publisher) fetchAndUpdatePrice() {
	newPrice := p.priceFetcher()
	p.lastPrice = p.price
	p.price = newPrice
	p.onPriceUpdated()
}
