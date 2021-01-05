package main

import (
	"testing"
	"time"

	"github.com/gen2brain/beeep"
)

func TestNotification(t *testing.T) {
	i := interval{5, 0, 3, 0, time.Now()}
	bannerText := "$33945.07 --> $33867.83 | Change: $-77.24 | Percent: -0.228%"
	hdr := sf("%d Minutes Passed | %.2f%%", i.MaxOccurences, i.PercentThreshold)
	beeep.Alert(hdr, bannerText, "assets/warning.png")
}
