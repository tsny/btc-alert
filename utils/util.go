package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	Up              = "🟩"
	Down            = "🟥"
	Neutral         = "🟦"
	dollar          = "💲"
	TimestampFormat = "03:04 PM"
)

func Banner(str string) {
	b := strings.Repeat("-", len(str))
	fmt.Printf("%s\n%s\n%s\n", b, str, b)
}

func Bannerf(str string, args ...interface{}) {
	str = fmt.Sprintf(str, args...)
	Banner(str)
}

func GetEmoji(curr, prev float64) string {
	if prev < curr {
		return Up
	} else if prev == curr {
		return Neutral
	}
	return Down
}

func GetTime() string {
	return time.Now().Format(TimestampFormat)
}
