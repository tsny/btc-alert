package main

import (
	"net/http"
	"os"

	"btc-alert/binance"
	"btc-alert/coinbase"
	"btc-alert/eps"
	"btc-alert/yahoo"

	log "github.com/sirupsen/logrus"
)

var lookupService = eps.NewSecurityLookup()

func main() {

	// Crypto
	for _, sec := range coinbase.CryptoMap {
		pub := eps.NewPublisher(coinbase.GetPrice, sec.Ticker, "Coinbase", false, 30)
		go trackSecurity(pub, sec)
	}

	// Stocks
	for _, ticker := range conf.YahooTickers {
		go func(ticker string) {
			pub := eps.NewPublisher(yahoo.GetPrice, ticker, "Yahoo", false, 30)
			pub.UseMarketHours = true
			name := yahoo.GetDetails(ticker).ShortName
			sec := eps.NewStock(name, ticker, "Coinbase")
			go trackSecurity(pub, sec)
		}(ticker)
	}

	// DOGE
	pub := eps.NewPublisher(binance.GetPrice, "DOGEUSDT", "Yahoo", false, 30)
	sec := eps.NewCrypto("DOGE", "DOGEUSDT", "Yahoo", "DOGE", "DOGECOIN")
	go trackSecurity(pub, sec)

	// API Init
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}

	trackGainers()
	log.Infof("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, initRoutes()))
}

func trackSecurity(pub *eps.Publisher, sec *eps.Security) *eps.InfoBall {
	_ = newListener(pub, conf.Intervals, conf.Thresholds)
	queue := eps.NewQueue().Subscribe(pub)
	info := lookupService.Register(sec, pub, queue)
	pub.SetActive(true)
	return info
}

// Finds the top gainers today from Yahoo and follows them
func trackGainers() {
	if !eps.IsMarketHours() {
		log.Warn("ignoring call to trackGainers as it is not market hours")
		return
	}
	gainers := yahoo.GetGainers()
	for _, v := range gainers {
		if info := lookupService.FindSecurityByNameOrTicker(v); info == nil {
			if deets := yahoo.GetDetails(v); deets != nil && deets.ShortName != "" {
				sec := eps.NewStock(deets.ShortName, v, "Yahoo")
				pub := eps.NewPublisher(yahoo.GetPrice, v, "Yahoo", true, 30)
				pub.UseMarketHours = true
				go trackSecurity(pub, sec)
			}
		}
	}
}
