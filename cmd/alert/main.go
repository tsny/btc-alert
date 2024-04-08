package main

import (
	"fmt"
	"os"
	"time"

	"btc-alert/pkg/alert"
	"btc-alert/pkg/coinbase"
	"btc-alert/pkg/eps"

	log "github.com/sirupsen/logrus"
)

var discordBot *alert.CryptoBot

func main() {

	if len(os.Args) > 1 {
		c := coinbase.Get24Hour(os.Args[1])
		fmt.Printf("%+v\n", c)
		return
	}
	conf := readConfig()

	btc := eps.NewPublisher(coinbase.GetPrice, coinbase.BTC, "Coinbase", false, 60, 20)
	alert.Publishers = append(alert.Publishers, btc)
	volListeners := []*alert.VolatilityListener{}
	for _, pc := range conf.PercentageChanges {
		listener := alert.NewVolatilityListener(btc, float64(pc.PercentChange), pc.DurInMinutes)
		volListeners = append(volListeners, listener)
	}
	btc.Start()

	go func() {
		time.Sleep(1 * time.Second)
		userID := conf.Discord.UsersToNotify[0]
		for {
			firstCandle := *btc.Candle
			dur := time.Hour * 6
			log.Infof("Alerting %v in %v", userID, dur)
			time.Sleep(dur)
			msg := firstCandle.DiffString(*btc.Candle)
			emoji := firstCandle.DiffEmoji(*btc.Candle)
			msg = fmt.Sprintf("%v: %v dur change: %v", emoji, dur, msg)
			_, err := discordBot.SendMessage(msg, userID)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}()

	for {
		select {}
	}

	// // API Init
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// 	log.Infof("Defaulting to port %s", port)
	// }

	// log.Infof("Listening on port %s", port)
	// log.Fatal(http.ListenAndServe(":"+port, initRoutes()))
}
