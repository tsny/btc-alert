package main

import (
	"net/http"
	"os"

	"btc-alert/binance"
	"btc-alert/coinbase"
	"btc-alert/eps"
	"btc-alert/priceTracking"
	"btc-alert/yahoo"

	log "github.com/sirupsen/logrus"
)

// PublisherMap is a WIP system for keeping track of all the
// EPS publishers for crypto/stockss
var PublisherMap = map[string]*eps.Publisher{}
var watchlist = map[string]*eps.Publisher{}
var queueService = priceTracking.NewQueueService()

func main() {
	// Crypto
	for _, ticker := range coinbase.CryptoMap {
		go trackSecurity(coinbase.GetPrice, ticker, "Coinbase", 10)
	}

	// Stocks
	for _, t := range conf.YahooTickers {
		go trackSecurity(yahoo.GetPrice, t, "Yahoo", 30)
	}

	go trackSecurity(binance.GetPrice, "DOGEUSDT", "Binance", 30)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}

	log.Infof("Listening on port %s", port)
	go http.ListenAndServe(":"+port, initRoutes())
	for {
	}
}

func trackSecurity(f func(string) float64, ticker, source string, dur int) {
	pub := eps.New(f, ticker, source, false, 30)
	_ = newListener(pub, conf.Intervals, conf.Thresholds)
	queueService.TrackSecurities(pub)
	PublisherMap[ticker] = pub
	pub.SetActive(true)
}

func refreshWatchlist() {
	tickers := yahoo.GetTopMoversTickers(true)
	for _, p := range watchlist {
		p.SetActive(false)
	}
	watchlist = make(map[string]*eps.Publisher)
	for _, t := range tickers {
		pub := eps.New(yahoo.GetPrice, t, "Yahoo", true, 30)
		pub.UseMarketHours = true
		_ = newListener(pub, conf.Intervals, conf.Thresholds)
		watchlist[t] = pub
	}
}
