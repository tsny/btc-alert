package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"btc-alert/coinbase"
	"btc-alert/eps"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
)

var lookupService = eps.NewSecurityLookup()
var publishers = []*eps.Publisher{}
var test *VolatilityListener

func main() {

	if len(os.Args) > 1 {
		c := coinbase.Get24Hour(os.Args[1])
		fmt.Printf("%+v\n", c)
		return
	}
	readConfig()

	// Crypto
	btc := eps.NewPublisher(coinbase.GetPrice, coinbase.BTC, "Coinbase", false, 60, 20)
	publishers = append(publishers, btc)
	for _, pc := range conf.PercentageChanges {
		test = NewVolatilityListener(btc, float64(pc.PercentChange), pc.DurInMinutes)
	}
	go track(btc)

	go func() {
		userID := conf.Discord.UsersToNotify[0]
		for {
			candle := btc.Candle
			dur := time.Hour * 6
			log.Infof("Alerting %v in %v", userID, dur)
			time.Sleep(dur)
			_, err := cryptoBot.SendMessage(btc.Candle.Diff(*candle), userID)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}()

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
	pub.Start()
}

func findPublisher(ticker string) (*eps.Publisher, bool) {
	ticker = strings.ToLower(ticker)
	return lo.Find(publishers, func(e *eps.Publisher) bool {
		return strings.Contains(strings.ToLower(e.Ticker), ticker)
	})
}
