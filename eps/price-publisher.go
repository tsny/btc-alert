// Package eps stands for Exchange Publisher Service
package eps

import (
	"fmt"
	"math"
	"time"

	log "github.com/sirupsen/logrus"
)

// Publisher periodically grabs data from its URL
// and sends out updates with the price it gets back
type Publisher struct {
	Ticker          string
	Source          string // Yahoo, Binance, etc
	UseMarketHours  bool   // whether the security abides by exchange market hours (NYSE, etc)
	CurrentCandle   *Candlestick
	Streak          int                             // How many times in a row the candlestick went up
	Updates         int                             // How many times the publisher has fetched a new price
	callbacks       []func(*Publisher, Candlestick) // callbacks upon price update
	closedCallbacks []func(*Publisher, Candlestick) // callbacks upon candlestick closed
	streakCallbacks []func(*Publisher, Candlestick, int)
	active          bool
	sleepDuration   int
	priceFetcher    func(string) float64
}

// NewPublisher is a constructor
func NewPublisher(priceFetcher func(string) float64, ticker, source string, start bool, sleepDur int) *Publisher {
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

func (p *Publisher) String() string {
	return fmt.Sprintf("%s [%s] - %d Updates - Active? [%v]", p.Ticker, p.Source, p.Updates, p.active)
}

// SetActive sets the publishers state
// active determines whether it will fetch and produce events
func (p *Publisher) SetActive(state bool) {
	p.active = state
}

// StartProducing loops and updates the price from the chosen exchange
func (p *Publisher) init() {
	firstRun := true
	go func() {
		p.active = true
		for {
			// Disable self if past market hours
			if p.UseMarketHours && !isMarketHours() && p.active {
				log.Warnf("%s disabled as it is not market hours", p.Ticker)
				p.active = false
			}
			if p.active {
				p.fetchAndUpdatePrice()
				if firstRun {
					curr := p.priceFetcher(p.Ticker)
					s := "%s -- Price Publisher [%s] active -- Current: %.2f\n"
					log.Infof(s, p.Ticker, p.Source, curr)
					firstRun = false
				}
			}
			time.Sleep(time.Duration(p.sleepDuration) * time.Second)
		}
	}()
}

// Regular US stock market trading hours are 9:30 AM -> 4 PM
func isMarketHours() bool {
	nyse, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(nyse)
	hour := now.Hour()
	min := now.Minute()
	return hour < 16 && (hour > 9 && min >= 30)
}

// Volatility returns the percent difference between the high/low and close
func (c Candlestick) Volatility() float64 {
	return (math.Abs(c.High-c.Low) / c.Close) * 100
}

// RegisterSubscriber assigns the func passed in to be called whenever
// the publisher has fetched and updated the price of the security
// todo: maybe this should take in an interface rather than just a func
func (p *Publisher) RegisterSubscriber(subscriber func(p *Publisher, c Candlestick)) {
	log.Infof("%s Publisher has new subscriber\n", p.Ticker)
	p.callbacks = append(p.callbacks, subscriber)
}

func (p *Publisher) onPriceUpdated() {
	p.Updates++
	for _, c := range p.callbacks {
		go c(p, *p.CurrentCandle)
	}
	if p.CurrentCandle.Complete {
		for _, c := range p.closedCallbacks {
			go c(p, *p.CurrentCandle)
		}
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
		log.Warnf("%s's price for %s was <= 0 \n", p.Source, p.Ticker)
		return
	}
	if p.CurrentCandle == nil {
		p.CurrentCandle = NewCandlestick(newPrice, p.sleepDuration, p.Ticker, p.Source)
	}
	candleDone := p.CurrentCandle.Update(newPrice)
	p.onPriceUpdated()
	if candleDone {
		if p.CurrentCandle.ClosedAboveOpen() {
			p.Streak++
		} else {
			p.Streak = 0
		}
		p.CurrentCandle = NewCandlestick(newPrice, p.sleepDuration, p.Ticker, p.Source)
	}
}
