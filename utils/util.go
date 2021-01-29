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

// Banner displays 'str' lined with '-'s
func Banner(str string) {
	b := strings.Repeat("-", len(str))
	fmt.Printf("%s\n%s\n%s\n", b, str, b)
}

// Bannerf is banner but with formatting args
func Bannerf(str string, args ...interface{}) {
	str = fmt.Sprintf(str, args...)
	Banner(str)
}

// Fts means float-to-string
func Fts(in float64) string {
	if in < 1 {
		return fmt.Sprintf("$%.4f", in)
	} else {
		return fmt.Sprintf("$%.2f", in)
	}
}

// GetEmoji gets a corresponding emoji based on price movement
func GetEmoji(curr, prev float64) string {
	if prev < curr {
		return Up
	} else if prev == curr {
		return Neutral
	}
	return Down
}

// GetTime gets the current time as string for logs
func GetTime() string {
	return time.Now().Format(TimestampFormat)
}
