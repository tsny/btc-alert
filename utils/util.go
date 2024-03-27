package utils

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

const (
	Up              = "↑"
	Down            = "↓"
	Neutral         = " "
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
	if math.Abs(in) < 1 {
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

func Getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		val = fallback
	}
	return val
}

func EnsureEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(key + " - env var missing")
	}
	return val
}

// Regular US stock market trading hours are 9:30 AM -> 4 PM
// TODO: Fix for 4->4:30, rn we just check for 9 to 4
func IsMarketHours() bool {
	nyse, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(nyse)
	if day := now.Weekday().String(); day == "Sunday" || day == "Saturday" {
		return false
	}
	hour := now.Hour()
	// min := now.Minute()
	return hour < 16 && hour > 9
}

func Fdate(t time.Time) string {
	return t.Format(time.DateTime)
}

func Ftime(t time.Time) string {
	return t.Format(time.TimeOnly)
}

func CompareTimes(t, t2 time.Time) string {
	return fmt.Sprintf("%v => %v", Ftime(t), Ftime(t2))
}
