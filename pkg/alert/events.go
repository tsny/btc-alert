package alert

import "time"

var bootTime = time.Now()

type PercentChange struct {
	DurInMinutes  int `json:"dur"`
	PercentChange int `json:"percentChange"`
}
