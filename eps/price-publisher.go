// Package eps stands for Exchange Publisher Service
package eps

import (
	"fmt"
	"math"
	"time"

	"github.com/tsny/btc-alert/utils"
)

// Publisher periodically grabs data from its URL
// and sends out updates with the price it gets back
type Publisher struct {
	Source          string // Yahoo, Binance, etc
	Ticker          string
	CurrentCandle   *Candlestick
	Streak          int // How many times in a row the candlestick went up
	callbacks       []func(*Publisher, Candlestick)
	streakCallbacks []func(*Publisher, Candlestick, int)
	active          bool
	sleepDuration   int
	priceFetcher    func(string) float64
}

// Candlestick represents a 'tick' or duration of a security's price
// https://en.wikipedia.org/wiki/Candlestick_chart
type Candlestick struct {
	Ticker            string
	DurationInSeconds int
	Begin             time.Time
	LastUpdate        time.Time
	Previous          float64
	Current           float64
	Close             float64
	Open              float64
	High              float64
	Low               float64
	Complete          bool
}

// NewCandlestick is a constructor
func NewCandlestick(open float64, dur int, ticker string) *Candlestick {
	return &Candlestick{
		Ticker:            ticker,
		DurationInSeconds: dur,
		Open:              open,
		Begin:             time.Now(),
	}
}

// Update checks the candlestick's lows/highs
// and returns whether or not the candlestick has completed
func (c *Candlestick) Update(price float64) bool {
	c.Previous = c.Current
	c.LastUpdate = time.Now()
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

// ClosedAboveOpen is whether or not the Candlestick
// completed with a price above it's open price
func (c Candlestick) ClosedAboveOpen() bool {
	return c.Close > c.Open
}

func (c Candlestick) String() string {
	emoji := utils.GetEmoji(c.Current, c.Previous)
	diff := c.Current - c.Previous
	percent := (diff / c.Current) * 100
	if c.Previous == 0.00 {
		return fmt.Sprintf("%s: (%s) $%.2f \n", emoji, c.Ticker, c.Current)
	}
	s := "%s %s: (%s) $%.2f | High: $%.2f | Low: $%.2f | Chg: $%.2f | Percent: %.2f%% | Volatility: %.2f%%"
	if c.Current < 1 {
		s = "%s: (%s) $%.5f | High: $%.5f | Low: $%.5f | Chg: $%.5f | Percent: %.2f%% | Volatility: %.2f%%"
	}
	return fmt.Sprintf(s, emoji, c.Ticker, c.Current, c.High, c.Low, diff, percent, c.Volatility())
}

// New is a constructor
func New(priceFetcher func(string) float64, ticker string, source string, start bool, sleepDur int) *Publisher {
	p := &Publisher{
		Source:        source,
		Ticker:        ticker,
		callbacks:     []func(p *Publisher, c Candlestick){},
		sleepDuration: sleepDur,
		priceFetcher:  priceFetcher,
		active:        start,
	}
	p.init()
	return p
}

// SetActive sets the publishers state
// active determines whether it will fetch and produce events
func (p *Publisher) SetActive(state bool) {
	p.active = state
}

// StartProducing loops and updates the price from the chosen exchange
func (p *Publisher) init() {
	curr := p.priceFetcher(p.Ticker)
	fmt.Printf("%s -- Price Publisher active -- Current: %.2f\n", p.Ticker, curr)
	go func() {
		for {
			if p.active {
				p.fetchAndUpdatePrice()
			}
			time.Sleep(time.Duration(p.sleepDuration) * time.Second)
		}
	}()
}

// Volatility returns the percent difference between the high/low and close
func (c Candlestick) Volatility() float64 {
	return (math.Abs(c.High-c.Low) / c.Close) * 100
}

// Subscribe assigns the func passed in to be called whenever
// the publisher has fetched and updated the price of the security
func (p *Publisher) Subscribe(f func(p *Publisher, c Candlestick)) {
	fmt.Printf("%s Publisher has new subscriber\n", p.Ticker)
	p.callbacks = append(p.callbacks, f)
}

func (p *Publisher) onPriceUpdated() {
	for _, c := range p.callbacks {
		go c(p, *p.CurrentCandle)
	}
}

// GetPrice returns the current price of the configured ticker
func (p *Publisher) GetPrice() float64 {
	if p.CurrentCandle == nil {
		return p.priceFetcher(p.Ticker)
	}
	return p.CurrentCandle.Current
}

func (p *Publisher) fetchAndUpdatePrice() {
	newPrice := p.priceFetcher(p.Ticker)
	// Ignore <= 0 since the API probably failed
	if newPrice <= 0 {
		fmt.Printf("warn: price for from %s for  %s was <= 0 \n", p.Source, p.Ticker)
		return
	}
	if p.CurrentCandle == nil {
		p.CurrentCandle = NewCandlestick(newPrice, p.sleepDuration, p.Ticker)
	}
	candleDone := p.CurrentCandle.Update(newPrice)
	p.onPriceUpdated()
	if candleDone {
		if p.CurrentCandle.ClosedAboveOpen() {
			p.Streak++
		} else {
			p.Streak = 0
		}
		p.CurrentCandle = NewCandlestick(newPrice, p.sleepDuration, p.Ticker)
	}
}
