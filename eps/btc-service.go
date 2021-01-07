// Package eps stands for Exchange Publisher Service
package eps

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tsny/btc-alert/coindesk"
)

const (
	// YahooURL = Yahoo Finance
	YahooURL = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=BTC-USD"
	// CoindeskURL = Coindesk API - Current Price endpoint
	CoindeskURL = "https://api.coindesk.com/v1/bpi/currentprice.json"
)

// Publisher periodically grabs data from its URL
// and sends out updates with the price it gets back
type Publisher struct {
	url             string
	priceUpdateSubs []func(new, old float64)
	price           float64
	lastPrice       float64
	active          bool
	sleepDuration   int
}

// New is a constructor
func New(url string) *Publisher {
	return &Publisher{
		url,
		[]func(new, old float64){},
		0,
		0,
		false,
		60,
	}
}

// StartListening loops and updates the price from the chosen exchange
func (p *Publisher) StartListening() {
	println("Exchange Price Publisher active")
	if p.active {
		return
	}
	for {
		p.fetchAndUpdatePrice()
		time.Sleep(time.Duration(p.sleepDuration) * time.Second)
	}
}

// Subscribe allows services to subscribe to new BitCoin events
func (p *Publisher) Subscribe(f func(new, old float64)) {
	println("Publisher has new subscriber")
	p.priceUpdateSubs = append(p.priceUpdateSubs, f)
}

func (p *Publisher) onPriceUpdated() {
	for _, c := range p.priceUpdateSubs {
		c(p.price, p.lastPrice)
	}
}

func (p *Publisher) fetchAndUpdatePrice() {
	res, err := http.Get(p.url)
	if err != nil {
		println(err)
		return
	}
	var out coindesk.Result
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil {
		panic(err)
	}
	s := strings.ReplaceAll(out.Bpi.USD.Rate, ",", "")
	p.lastPrice = p.price
	p.price, err = strconv.ParseFloat(s, 64)
	if err != nil {
		println(err.Error())
		return
	}
	p.onPriceUpdated()
}
