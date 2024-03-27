// Package eps stands for Exchange Publisher Service
package eps

import (
	"btc-alert/utils"
	"fmt"
	"strings"
	"time"

	"github.com/blend/go-sdk/uuid"
	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

// Publisher periodically grabs data from its URL
// and sends out updates with the price it gets back
type Publisher struct {
	Ticker         string
	Source         string // Yahoo, Binance, etc
	UseMarketHours bool   // whether the security abides by exchange market hours (NYSE, etc)
	Candle         *Candlestick
	PreviousCandle *Candlestick
	Listeners      map[string]UpdateHandler
	Stack          *CandleStack

	// Rethink:
	Streak         int // How many times in a row the candlestick moved in a certain direction
	PositiveStreak bool

	active       bool
	refreshDur   time.Duration
	candleDur    time.Duration
	priceFetcher func(string) float64
}

type UpdateHandler func(*Publisher, *Candlestick, bool)

// NewPublisher is a constructor
func NewPublisher(priceFetcher func(string) float64, ticker,
	source string, start bool, candleDurSeconds int, refreshSeconds int) *Publisher {
	if candleDurSeconds < 0 {
		candleDurSeconds = 60
	}
	if refreshSeconds < 0 {
		refreshSeconds = 15
	}
	p := &Publisher{
		Source:       source,
		Ticker:       strings.ToLower(ticker),
		refreshDur:   time.Duration(refreshSeconds) * time.Second,
		candleDur:    time.Duration(candleDurSeconds) * time.Second,
		priceFetcher: priceFetcher,
		active:       start,
		Listeners:    make(map[string]UpdateHandler),
	}
	p.Stack = NewStack(p)
	log.Infof("Publisher %v: %v %v", p.Ticker, p.refreshDur, p.candleDur)
	if start {
		p.Start()
	}
	return p
}

func (p *Publisher) GetRefreshDurInSeconds() int {
	return int(p.refreshDur.Seconds())
}

func (p *Publisher) String() string {
	candleString := ""
	if p.Candle == nil {
		candleString = "no candle yet"
	} else {
		candleString = p.Candle.String()
	}
	s := "%v [%v] (%v) - Active? [%v]"
	return fmt.Sprintf(s, p.Ticker, p.Source, candleString, p.active)
}

func (p *Publisher) Unsub(id string) bool {
	if lo.Contains(lo.Keys(p.Listeners), id) {
		delete(p.Listeners, id)
		return true
	}
	return false
}

// SetActive sets the publishers state
// active determines whether it will fetch and produce events
func (p *Publisher) SetActive(state bool) {
	p.active = state
}

// StartProducing loops and updates the price from the chosen exchange
func (p *Publisher) Start() {
	log.Infof("%v publisher: starting", p.Ticker)
	p.active = true
	go func() {
		for {
			// Disable self if past market hours
			if p.UseMarketHours && !utils.IsMarketHours() && p.active {
				log.Warnf("%s disabled as it is not market hours", p.Ticker)
				p.active = false
			}
			done := false
			if p.active {
				done = p.fetchAndUpdatePrice()
			}
			s := fmt.Sprintf("%v: update [%v] %v => %v",
				p.Ticker, p.Price(), fdate(p.Candle.Start), fdate(time.Now().Local().Add(p.refreshDur)))
			if done {
				s += " [done]"
			}
			log.Infof(s)
			time.Sleep(p.refreshDur)
		}
	}()
}

func fdate(t time.Time) string {
	return t.Format(time.DateTime)
}

func (p *Publisher) RegisterPriceUpdateListener(s UpdateHandler) string {
	handlerID := uuid.V4().String()
	p.Listeners[handlerID] = s
	return handlerID
}

func (p *Publisher) onPriceUpdated(completed bool) {
	for _, c := range p.Listeners {
		go c(p, p.Candle, completed)
	}
}

// Price returns the current price of the configured ticker
func (p *Publisher) Price() float64 {
	if p.Candle == nil {
		return p.priceFetcher(p.Ticker)
	}
	return p.Candle.Price
}

func (p *Publisher) fetchAndUpdatePrice() bool {
	newPrice := p.priceFetcher(p.Ticker)
	// Ignore <= 0 since the API probably failed
	if newPrice <= 0 {
		log.Errorf("%s's price for %s was <= 0", p.Source, p.Ticker)
		return false
	}
	if p.Candle == nil {
		p.newCandlestick(newPrice)
		return false
	}
	p.Candle.Update(newPrice)
	candleDone := p.Candle.Start.Add(p.candleDur).Before(time.Now().Local())
	if candleDone {
		p.newCandlestick(newPrice)
	}
	p.onPriceUpdated(candleDone)
	return candleDone
}

func (p *Publisher) newCandlestick(price float64) {
	p.checkStreak()
	p.PreviousCandle = p.Candle
	if p.PreviousCandle != nil {
		p.PreviousCandle.Finish()
	}
	p.Candle = NewCandlestick(price)
}

// Checks whether a security is streaking in either direction
func (p *Publisher) checkStreak() {
	if p.Candle == nil {
		return
	}
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
