package main

import (
	"fmt"
	. "fmt"
	"math"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

var sf = Sprintf

func onPriceUpdated() {
	intervalCompleted := false
	for _, i := range conf.Intervals {
		if i.beginPrice == 0 {
			i.beginPrice = price
		}
		i.occurrences++
		if i.occurrences >= i.MaxOccurences {
			i.onCompleted()
			i.reset()
			intervalCompleted = true
		}
	}
	for _, c := range conf.Thresholds {
		c.checkThreshold()
	}
	if !intervalCompleted {
		fmt.Print(getSummary())
		if conf.UseDiscord {
			discordMessage(getSummary(), false)
		}
	}
}

func getSummary() string {
	t := getTime()
	emoji := getEmoji(price, lastPrice)
	diff := price - lastPrice
	percent := (diff / price) * 100
	if lastPrice == 0.00 {
		return sf("%s %s: $%.2f \n", emoji, t, price)
	}
	return sf("%s %s: $%.2f | Change: $%.2f | Percent: %.3f%% \n", emoji, t, price, diff, percent)
}

func (i *interval) onCompleted() {
	diff := price - i.beginPrice
	percent := (diff / i.beginPrice) * 100
	prefix := ""
	if math.Abs(percent) > i.PercentThreshold {
		prefix = sf("%s ALERT: Threshold of %.2f%% was reached! ", alert, i.PercentThreshold)
	}

	totalChange := sf("$%.2f --> $%.2f", i.beginPrice, price)
	changes := sf("Change: $%.2f | Percent: %.3f%%", diff, percent)
	bannerText := sf("%s: %s%d Minutes Passed | %s | %s", getTime(), prefix, i.occurrences, totalChange, changes)
	banner(bannerText)

	if math.Abs(percent) > i.PercentThreshold {
		hdr := sf("%d Minutes Passed | %.2f%%", i.MaxOccurences, i.PercentThreshold)
		notif(hdr, bannerText, "assets/warning.png")
		if conf.UseDiscord {
			discordMessage(hdr, true)
			discordMessage(bannerText, false)
		}
	} else {
		discordMessage(bannerText, false)
	}
}

func (t *threshold) checkThreshold() {
	if price >= t.Threshold+t.beginPrice {
		t.onThresholdReached(true)
	} else if price <= t.beginPrice-t.Threshold {
		t.onThresholdReached(false)
	}
}

func (t *threshold) onThresholdReached(breachedUp bool) {
	emoji := down
	if breachedUp {
		emoji = up
	}

	hdr := sf("Price Threshold Breached: $%v", t.Threshold)
	str := "%s %s: %s | %s ($%.2f)"
	body := sf(str, emoji, getTime(), hdr, fpm(t.beginPrice, price), t.beginPrice-price)

	notif(hdr, body, "assets/warning.png")

	banner("ALERT " + body)

	if conf.UseDiscord {
		discordMessage(body, false)
	}
	t.beginPrice = price
}

// fpm -- formatPriceMovement
func fpm(begin, end float64) string {
	return sf("$%.2f --> $%.2f", begin, end)
}

func getTime() string {
	return time.Now().Format(format)
}

func (i *interval) reset() {
	i.occurrences = 0
	i.startTime = time.Now()
	i.beginPrice = price
}

// notif is a beeep.Alert() wrapper
// it ensures there are no '$'
// since this can mess with powershell notifications
func notif(hdr, body, img string) {
	hdr = strings.ReplaceAll(hdr, "$", "💲")
	body = strings.ReplaceAll(body, "$", "💲")
	body = strings.ReplaceAll(body, "|", "\n")
	beeep.Alert(hdr, body, img)
}
