package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/tsny/btc-alert/eps"
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
		prefix = sf("%s ALERT: %.2f%%! ", alert, i.PercentThreshold)
	}

	totalChange := sf("$%.2f --> $%.2f", i.beginPrice, new)
	changes := sf("Chg: $%.2f | Percent: %.3f%%", diff, percent)

	bannerText := sf("%s: (%s) %s%d Min Passed | %s | %s",
		getTime(), p.Source, prefix, i.occurrences, totalChange, changes)
	banner(bannerText)

	if math.Abs(percent) > i.PercentThreshold {
		if conf.DesktopNotifications {
			hdr := sf("Last %d Min | %.2f%%", i.MaxOccurences, i.PercentThreshold)
			notif(hdr, bannerText, "assets/warning.png")
		}
		if conf.Discord.Enabled {
			discordMessage(bannerText, true)
		}
	} else {
		if conf.Discord.Enabled {
			discordMessage(bannerText, false)
		}
	}
}

func (t *threshold) onThresholdReached(p *eps.Publisher, breachedUp bool, new, old float64) {
	emoji := down
	if breachedUp {
		emoji = up
	}

	hdr := sf("Price Movement: $%v", t.Threshold)
	str := "%s %s: (%s) %s | %s ($%.2f)"
	body := sf(str, emoji, getTime(), p.Source, hdr, fpm(t.beginPrice, new), new-t.beginPrice)

	if conf.DesktopNotifications {
		notif(hdr, body, "assets/warning.png")
	}
	banner("ALERT " + body)

	if conf.Discord.Enabled {
		discordMessage(body, false)
	}
	t.beginPrice = new
}

// fpm -- formatPriceMovement
func fpm(begin, end float64) string {
	return sf("$%.2f --> $%.2f", begin, end)
}

func getTime() string {
	return time.Now().Format(format)
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
