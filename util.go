package main

import (
	"fmt"
	"strings"
)

const (
	up      = "🟩"
	down    = "🟥"
	neutral = "🟦"
	alert   = "☎️"
	dollar  = "💲"
	format  = "03:04 PM"
)

func banner(str string) {
	b := strings.Repeat("-", len(str))
	fmt.Printf("%s\n%s\n%s\n", b, str, b)
}

func bannerf(str string, args ...interface{}) {
	str = sf(str, args...)
	banner(str)
}

func getEmoji(curr, prev float64) string {
	if prev < curr {
		return up
	} else if prev == curr {
		return neutral
	}
	return down
}
