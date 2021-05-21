package eps

import (
	"btc-alert/utils"
	"fmt"
	"math"
	"time"
)

// Candlestick represents a 'tick' or duration of a security's price
// https://en.wikipedia.org/wiki/Candlestick_chart
type Candlestick struct {
	Ticker            string
	Source            string
	DurationInSeconds int
	Begin             time.Time
	LastUpdate        time.Time
	Previous          float64
	Price             float64
	Close             float64
	Open              float64
	High              float64
	Low               float64
	Complete          bool
}

// NewCandlestick is a constructor
func NewCandlestick(open float64, dur int, ticker, source string) *Candlestick {
	return &Candlestick{
		Ticker:            ticker,
		Source:            source,
		DurationInSeconds: dur,
		Open:              open,
		Begin:             time.Now(),
	}
}

// Update checks the candlestick's lows/highs
// and returns whether or not the candlestick has completed
func (c *Candlestick) Update(price float64) bool {
	c.Previous = c.Price
	c.LastUpdate = time.Now()
	c.Price = price
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

// OpenCloseDiff returns the absolute difference between the candle's
// close and open
func (c Candlestick) OpenCloseDiff() float64 {
	return math.Abs(c.Close - c.Open)
}

func (c Candlestick) String() string {
	emoji := utils.GetEmoji(c.Price, c.Previous)
	diff := c.Price - c.Previous
	percent := (diff / c.Price) * 100
	if c.Previous == 0.00 {
		return fmt.Sprintf("%s: (%s) $%.2f \n", emoji, c.Ticker, c.Price)
	}
	s := "%s (%s) $%.2f | High: $%.2f | Low: $%.2f | Chg: $%.2f | Percent: %.2f%% | Vol: %.2f%%"
	// Larger formatting for crypto prices below $1
	if c.Price < 1 {
		s = "%s (%s) $%.5f | High: $%.5f | Low: $%.5f | Chg: $%.5f | Percent: %.2f%% | Vol: %.2f%%"
	}
	return fmt.Sprintf(s, emoji, c.Ticker, c.Price, c.High, c.Low, diff, percent, c.Volatility())
}
