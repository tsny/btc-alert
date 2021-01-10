package main

import (
	"testing"
	"time"
)

// func TestThreshold(t *testing.T) {
// 	c := threshold{544, 400}
// 	c.onThresholdReached(true, 300, 499)
// }

func TestNotification(t *testing.T) {
	i := interval{5, 0, 3, 0, time.Now()}
	bannerText := "$33945.07 --> $33867.83 | Change: $-77.24 | Percent: -0.228%"
	hdr := sf("%d Minutes Passed | %.2f%%", i.MaxOccurences, i.PercentThreshold)
	notif(hdr, bannerText, "assets/warning.png")
}
