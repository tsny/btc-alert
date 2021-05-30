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
	Candle          *Candlestick
	PreviousCandle  *Candlestick
	Streak          int // How many times in a row the candlestick moved in a certain direction
	PositiveStreak  bool
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
	var price float64
	if p.Candle == nil {
		price = p.GetPrice()
	} else {
		price = p.Candle.Price
	}
	s := "%s [%s] (%v) - %d Updates - Active? [%v]"
	return fmt.Sprintf(s, p.Ticker, p.Source, price, p.Updates, p.active)
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
			if p.UseMarketHours && !IsMarketHours() && p.active {
				log.Warnf("%s disabled as it is not market hours", p.Ticker)
				p.active = false
			}
			if p.active {
				p.fetchAndUpdatePrice()
				if firstRun {
					curr := p.priceFetcher(p.Ticker)
					s := "%s Price Publisher [%s] active -- Current: %.2f\n"
					log.Infof(s, p.Ticker, p.Source, curr)
					firstRun = false
				}
			}
			time.Sleep(time.Duration(p.sleepDuration) * time.Second)
		}
	}()
}

// Regular US stock market trading hours are 9:30 AM -> 4 PM
// TODO: Fix for 4->4:30, rn we just check for 9 to 4
func IsMarketHours() bool {
	nyse, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(nyse)
	if day := now.Weekday().String(); day == "Sunday" || day == "Saturday" {
		return false
	}
	hour := now.Hour()
	// min := now.Minute()
	return hour < 16 && hour > 9
}

// Volatility returns the percent difference between the high/low and close
func (c Candlestick) Volatility() float64 {
	return (math.Abs(c.High-c.Low) / c.Close) * 100
}

// RegisterSubscriber assigns the func passed in to be called whenever
// the publisher has fetched and updated the price of the security
// todo: maybe this should take in an interface rather than just a func
func (p *Publisher) RegisterSubscriber(subscriber func(p *Publisher, c Candlestick)) {
	p.callbacks = append(p.callbacks, subscriber)
}

func (p *Publisher) onPriceUpdated() {
	p.Updates++
	for _, c := range p.callbacks {
		go c(p, *p.Candle)
	}
	if p.Candle.Complete {
		for _, c := range p.closedCallbacks {
			go c(p, *p.Candle)
		}
	}
}

// GetPrice returns the current price of the configured ticker
func (p *Publisher) GetPrice() float64 {
	if p.Candle == nil {
		return p.priceFetcher(p.Ticker)
	}
	return p.Candle.Price
}

func (p *Publisher) fetchAndUpdatePrice() {
	newPrice := p.priceFetcher(p.Ticker)
	// Ignore <= 0 since the API probably failed
	if newPrice <= 0 {
		log.Warnf("%s's price for %s was <= 0 \n", p.Source, p.Ticker)
		return
	}
	if p.Candle == nil {
		p.Candle = NewCandlestick(newPrice, p.sleepDuration, p.Ticker, p.Source)
	}
	candleDone := p.Candle.Update(newPrice)
	p.onPriceUpdated()
	if candleDone {
		p.checkStreak()
		p.PreviousCandle = p.Candle
		p.Candle = NewCandlestick(newPrice, p.sleepDuration, p.Ticker, p.Source)
	}
}

// Checks whether a security is streaking in either direction
func (p *Publisher) checkStreak() {
	// If there is no streak, begin it
	if p.Streak == 0 {
		p.Streak++
		p.PositiveStreak = p.Candle.ClosedAboveOpen()
		return
	}
	// Otherwise, make sure the streak is still going in whichever direction
	if p.Candle.ClosedAboveOpen() && p.PositiveStreak {
		p.Streak++
	} else if !p.Candle.ClosedAboveOpen() && !p.PositiveStreak {
		p.Streak++
	} else {
		p.Streak = 0
	}
}

// Returns a summary string of the security's current streak
// Ex: BTC-USD (45023.23) [Coinbase] has a streak of -8
func (p *Publisher) StreakSummary() string {
	status := "-"
	if p.PositiveStreak {
		status = "+"
	}
	str := "%s (%.2f) [%s] has a streak of %s%d"
	return fmt.Sprintf(str, p.Ticker, p.Candle.Price, p.Source, status, p.Streak)
}
