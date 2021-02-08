package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/utils"
)

var sf = fmt.Sprintf

// thresholds are price jumps
// they alert after prices move a certain amount from
// the starting price
type threshold struct {
	beginPrice float64
	Threshold  float64 `json:"threshold"`
}

// intervals are checked every minute
// if 'maxOccurences' number of minutes pass
// then the interval lapses and onCompleted() is called
type interval struct {
	beginPrice       float64
	low              float64
	high             float64
	occurrences      int
	MaxOccurences    int     `json:"maxOccurences"`
	PercentThreshold float64 `json:"percentThreshold"`
	startTime        time.Time
}

func (i *interval) onCompleted(p *eps.Publisher, new, old float64) {
	diff := new - i.beginPrice
	percent := (diff / i.beginPrice) * 100
	prefix := ""
	if math.Abs(percent) > i.PercentThreshold {
		prefix = "ALERT"
	}

	totalChange := sf("%s --> %s", utils.Fts(i.beginPrice), utils.Fts(new))
	changes := sf("Chg: %s | Percent: %.3f%%", utils.Fts(diff), percent)

	bannerText := sf("%s: (%s) %s%d Min | %s | %s",
		utils.GetTime(), p.Ticker, prefix, i.occurrences, totalChange, changes)
	utils.Banner(bannerText)

	if math.Abs(percent) > i.PercentThreshold {
		if conf.DesktopNotifications {
			hdr := sf("Last %d Min | %.2f%%", i.MaxOccurences, i.PercentThreshold)
			notif(hdr, bannerText, "assets/warning.png")
		}
		if conf.Discord.Enabled {
			cryptoBot.SendMessage(bannerText, "everyone", false)
		}
	} else {
		if conf.Discord.Enabled {
			cryptoBot.SendMessage(bannerText, "", false)
		}
	}
}

func (t *threshold) onThresholdReached(p *eps.Publisher, breachedUp bool, new, old float64) {
	emoji := utils.Down
	if breachedUp {
		emoji = utils.Up
	}

	priceMovement := sf("Price Movement: $%v", t.Threshold)
	str := "%s %s: (%s) %s | %s ($%.2f)"
	body := sf(str, emoji, utils.GetTime(), p.Ticker, priceMovement, fpm(t.beginPrice, new), new-t.beginPrice)

	if conf.DesktopNotifications {
		notif(priceMovement, body, "assets/warning.png")
	}
	utils.Banner("ALERT " + body)

	if conf.Discord.Enabled {
		cryptoBot.SendMessage(body, "", false)
	}
	t.beginPrice = new
}

// fpm -- formatPriceMovement
func fpm(begin, end float64) string {
	return sf("$%.2f --> $%.2f", begin, end)
}

func (i *interval) reset(new float64) {
	i.occurrences = 0
	i.startTime = time.Now()
	i.beginPrice = new
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
