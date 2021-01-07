package coindesk

import "time"

type Result struct {
	Time       Time   `json:"time"`
	Disclaimer string `json:"disclaimer"`
	ChartName  string `json:"chartName"`
	Bpi        Bpi    `json:"bpi"`
}

type Time struct {
	Updated    string    `json:"updated"`
	UpdatedISO time.Time `json:"updatedISO"`
	Updateduk  string    `json:"updateduk"`
}

type USD struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

type GBP struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

type EUR struct {
	Code        string  `json:"code"`
	Symbol      string  `json:"symbol"`
	Rate        string  `json:"rate"`
	Description string  `json:"description"`
	RateFloat   float64 `json:"rate_float"`
}

type Bpi struct {
	USD USD `json:"USD"`
	GBP GBP `json:"GBP"`
	EUR EUR `json:"EUR"`
}
