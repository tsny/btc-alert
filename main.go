package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"btc-alert/coinbase"
	"btc-alert/eps"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

var lookupService = eps.NewSecurityLookup()
var publishers = []*eps.Publisher{}
var queues = map[string]*eps.CandleQueue{}

func main() {

	if len(os.Args) > 1 {
		c := coinbase.Get24Hour(os.Args[1])
		fmt.Printf("%+v\n", c)
		return
	}
	readConfig()

	// Crypto
	btcPublisher := eps.NewPublisher(coinbase.GetPrice, coinbase.BTC, "Coinbase", false, 60, 20)
	publishers = append(publishers, btcPublisher)
	go track(btcPublisher)

	// API Init
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Infof("Defaulting to port %s", port)
	}

	log.Infof("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, initRoutes()))
}

func track(pub *eps.Publisher) {
	// _ = newListener(pub, conf.Intervals, conf.Thresholds)
	queue := eps.NewQueue(pub)
	queues[pub.Ticker] = queue
	pub.SetActive(true)
}

func findPublisher(s string) (*eps.Publisher, bool) {
	return lo.Find(publishers, func(e *eps.Publisher) bool {
		return strings.Contains(strings.ToLower(e.Ticker), s)
	})
}

func findQueue(ticker string) (*eps.CandleQueue, bool) {
	key, ok := lo.Find(lo.Keys(queues), func(s string) bool { return strings.Contains(s, ticker) })
	if !ok {
		return nil, false
	}
	return queues[key], true
}
