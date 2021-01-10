// Package eps stands for Exchange Publisher Service
package eps

import (
	"fmt"
	"time"

	"github.com/tsny/btc-alert/utils"
)

// Publisher periodically grabs data from its URL
// and sends out updates with the price it gets back
type Publisher struct {
	Source        string
	callbacks     []func(*Publisher, Candlestick)
	active        bool
	sleepDuration int
	CurrentCandle *Candlestick
	priceFetcher  func() float64
}

// Candlestick represents a 'tick' or duration of a security's price
// https://en.wikipedia.org/wiki/Candlestick_chart
type Candlestick struct {
	Source            string
	DurationInSeconds int
	Begin             time.Time
	Previous          float64
	Current           float64
	Close             float64
	Open              float64
	High              float64
	Low               float64
	Complete          bool
}

// NewCandlestick is a constructor
func NewCandlestick(open float64, dur int, source string) *Candlestick {
	return &Candlestick{
		Source:            source,
		DurationInSeconds: dur,
		Open:              open,
		Begin:             time.Now(),
	}
}

// Update checks the candlestick's lows/highs
// and returns whether or not the candlestick has completed
func (c *Candlestick) Update(price float64) bool {
	c.Previous = c.Current
	c.Current = price
	if price < c.Low || c.Low == 0 {
		c.Low = price
	}
	if price > c.High {
		c.High = price
	}
	if time.Since(c.Begin).Seconds() >= 60 {
		c.Close = price
		c.Complete = true
		return true
	}
	return false
}

// ClosedAbove is whether or not the Candlestick
// completed with a price above it's open price
func (c Candlestick) ClosedAbove() bool {
	return c.Close > c.Open
}

func (c Candlestick) String() string {
	now := utils.GetTime()
	emoji := utils.GetEmoji(c.Current, c.Previous)
	diff := c.Current - c.Previous
	percent := (diff / c.Current) * 100
	if c.Previous == 0.00 {
		return fmt.Sprintf("%s %s: (%s) $%.2f \n", emoji, now, c.Source, c.Current)
	}
	s := "%s %s: (%s) $%.2f | High: $%.2f | Low: $%.2f | Chg: $%.2f | Percent: %.3f%% \n"
	return fmt.Sprintf(s, emoji, now, c.Source, c.Current, c.High, c.Low, diff, percent)
}

// New is a constructor
func New(priceFetcher func() float64, source string) *Publisher {
	return &Publisher{
		Source:        source,
		callbacks:     []func(p *Publisher, c Candlestick){},
		sleepDuration: 5,
		priceFetcher:  priceFetcher,
	}
}

// StartProducing loops and updates the price from the chosen exchange
func (p *Publisher) StartProducing() {
	fmt.Printf("%s -- Price Event Publisher active\n", p.Source)
	if p.active {
		fmt.Printf("%s -- Price Event Publisher is ALREADY active\n", p.Source)
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
func (p *Publisher) Subscribe(f func(p *Publisher, c Candlestick)) {
	fmt.Printf("%s Publisher has new subscriber\n", p.Source)
	p.callbacks = append(p.callbacks, f)
}

func (p *Publisher) onPriceUpdated() {
	for _, c := range p.callbacks {
		c(p, *p.CurrentCandle)
	}
}

func (p *Publisher) fetchAndUpdatePrice() {
	newPrice := p.priceFetcher()
	if p.CurrentCandle == nil {
		p.CurrentCandle = NewCandlestick(newPrice, p.sleepDuration, p.Source)
	}
	candleDone := p.CurrentCandle.Update(newPrice)
	p.onPriceUpdated()
	if candleDone {
		p.CurrentCandle = NewCandlestick(newPrice, p.sleepDuration, p.Source)
	}
}
