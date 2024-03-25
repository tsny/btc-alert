package eps

import (
	"btc-alert/utils"
	"fmt"
	"math"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

// Candlestick represents a 'tick' or duration of a security's price
// This project is primarily concerned with minute candles
// https://en.wikipedia.org/wiki/Candlestick_chart
type Candlestick struct {
	Start      time.Time
	End        time.Time
	LastUpdate time.Time
	Previous   float64
	Price      float64
	Close      float64
	Open       float64
	High       float64
	Low        float64
}

// NewCandlestick is a constructor
func NewCandlestick(open float64) *Candlestick {
	return &Candlestick{
		Open:  open,
		Price: open,
		Low:   open,
		High:  open,
		Start: time.Now().Local(),
	}
}

func (c *Candlestick) IsComplete() bool {
	return c.Close != 0
}

func (c *Candlestick) Finish() {
	c.End = time.Now().Local()
}

// Update checks the candlestick's lows/highs
// and returns whether or not the candlestick has completed
func (c *Candlestick) Update(price float64) {
	c.Previous = c.Price
	c.LastUpdate = time.Now()
	c.Price = price
	if price < c.Low || c.Low == 0 {
		c.Low = price
	}
	if price > c.High {
		c.High = price
	}
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
	if c.Previous == 0.00 && c.Price == 0.00 {
		return fmt.Sprintf("%v $%.2f", emoji, c.Price)
	}
	s := "%s $%.2f | High: $%.2f | Low: $%.2f | Chg: $%.2f | Percent: %.2f%% | Vol: %v"
	// Larger formatting for crypto prices below $1
	if c.Price < 1 {
		s = "%v $%.5f | High: $%.5f | Low: $%.5f | Chg: $%.5f | Percent: %.2f%% | Vol: %v%"
	}
	return fmt.Sprintf(s, emoji, c.Price, c.High, c.Low, diff, percent, c.Volatility())
}

func (c Candlestick) Table() string {
	t := table.NewWriter()
	t.AppendRow(table.Row{c.Price})
	return t.Render()
}

// Volatility returns the percent difference between the high/low and close
func (c Candlestick) Volatility() float64 {
	return (math.Abs(c.High-c.Low) / c.Close) * 100
}

func (c Candlestick) Diff(c2 *Candlestick) string {
	s := fmt.Sprintf("%v => %v (%v) | %v => %v",
		fdate(c.Start), fdate(c2.Start), c2.Start.Sub(c.Start), c.Open, c2.Open)
	return s
}
