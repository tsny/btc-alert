package main

import (
	"btc-alert/priceTracking"
	"log"
	"net/http"
	"os"

	"github.com/tsny/btc-alert/binance"
	"github.com/tsny/btc-alert/coinbase"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/utils"
	"github.com/tsny/btc-alert/yahoo"
)

// PublisherMap is a WIP system for keeping track of all the
// EPS publishers for crypto/stockss
var PublisherMap = map[string]*eps.Publisher{}
var watchlist = map[string]*eps.Publisher{}

func main() {
	if conf.DesktopNotifications {
		// notif("BTC-ALERT", "Desktop Notifications Enabled", "")
	}
	queueService := priceTracking.NewQueueService()
	test := eps.Publisher{}
	queueService.TrackSecurities(&test)

	// Crypto
	for _, ticker := range coinbase.CryptoMap {
		pub := eps.New(coinbase.GetPrice, ticker, "Coinbase", true, 10)
		_ = newListener(pub, conf.Intervals, conf.Thresholds)
		PublisherMap[ticker] = pub
	}

	// Stocks
	for _, t := range conf.YahooTickers {
		pub := eps.New(yahoo.GetPrice, t, "Yahoo", true, 30)
		pub.UseMarketHours = true
		_ = newListener(pub, conf.Intervals, conf.Thresholds)
		PublisherMap[t] = pub
	}

	// Doge
	pub := eps.New(binance.GetPrice, "DOGEUSDT", "Binance", true, 30)
	_ = newListener(pub, conf.Intervals, conf.Thresholds)
	PublisherMap["DOGE"] = pub

	utils.Banner("btc-alert initialized")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	go http.ListenAndServe(":"+port, initRoutes())
	for {
	}
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
