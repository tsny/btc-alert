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
// EPS publishers for crypto/stocks/
// TODO: It should most likely be a service that you can call with a ticker
// or name and it finds the security for you via lookup
var PublisherMap = map[string]*eps.Publisher{}
var lookupService = eps.NewSecurityLookup()

// TODO: probably redo this stuff too
var watchlist = map[string]*eps.Publisher{}
var queueService = priceTracking.NewQueueService()

func main() {

	// Crypto
	for _, sec := range coinbase.CryptoMap {
		pub := eps.NewPublisher(coinbase.GetPrice, sec.Ticker, "Coinbase", false, 30)
		// sec := eps.NewCrypto(name, ticker, "Coinbase", strings.ReplaceAll(ticker, "-USD", ""))
		lookupService.Register(sec, pub)
		go trackSecurity(pub)
	}

	// Stocks
	for _, ticker := range conf.YahooTickers {
		go func(ticker string) {
			pub := eps.NewPublisher(yahoo.GetPrice, ticker, "Yahoo", false, 30)
			pub.UseMarketHours = true
			name := yahoo.GetDetails(ticker).ShortName
			sec := eps.NewCrypto(name, ticker, "Coinbase")
			lookupService.Register(sec, pub)
			go trackSecurity(pub)
		}(ticker)
	}

	// DOGE
	pub := eps.NewPublisher(binance.GetPrice, "DOGEUSDT", "Yahoo", false, 30)
	sec := eps.NewCrypto("DOGE", "DOGEUSDT", "Yahoo", "DOGE", "DOGECOIN")
	lookupService.Register(sec, pub)
	go trackSecurity(pub)

	// API Init
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

func trackSecurity(pub *eps.Publisher) {
	_ = newListener(pub, conf.Intervals, conf.Thresholds)
	queueService.TrackSecurities(pub)
	PublisherMap[pub.Ticker] = pub
	pub.SetActive(true)
}

func refreshWatchlist() {
	// tickers := yahoo.GetTopMoversTickers(true)
	// for _, p := range watchlist {
	// 	p.SetActive(false)
	// }
	// watchlist = make(map[string]*eps.Publisher)
	// for _, t := range tickers {
	// pub := eps.New(yahoo.GetPrice, t, "Yahoo", true, 30)
	// pub.UseMarketHours = true
	// _ = newListener(pub, conf.Intervals, conf.Thresholds)
	// watchlist[t] = pub
	// }
}
