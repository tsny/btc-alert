package yahoo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// YahooURL = Yahoo Finance
	YahooURL = "https://query1.finance.yahoo.com/v7/finance/quote?symbols=%s"
)

// TLR is the Top level json result
type TLR struct {
	QuoteResponse QuoteResponse `json:"quoteResponse"`
}

type Result struct {
	Language                          string  `json:"language"`
	Region                            string  `json:"region"`
	QuoteType                         string  `json:"quoteType"`
	QuoteSourceName                   string  `json:"quoteSourceName"`
	Triggerable                       bool    `json:"triggerable"`
	Currency                          string  `json:"currency"`
	Exchange                          string  `json:"exchange"`
	ShortName                         string  `json:"shortName"`
	MessageBoardID                    string  `json:"messageBoardId"`
	ExchangeTimezoneName              string  `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName         string  `json:"exchangeTimezoneShortName"`
	GmtOffSetMilliseconds             int     `json:"gmtOffSetMilliseconds"`
	Market                            string  `json:"market"`
	EsgPopulated                      bool    `json:"esgPopulated"`
	FirstTradeDateMilliseconds        int64   `json:"firstTradeDateMilliseconds"`
	PriceHint                         int     `json:"priceHint"`
	CirculatingSupply                 int     `json:"circulatingSupply"`
	LastMarket                        string  `json:"lastMarket"`
	Volume24Hr                        int64   `json:"volume24Hr"`
	VolumeAllCurrencies               int64   `json:"volumeAllCurrencies"`
	FromCurrency                      string  `json:"fromCurrency"`
	ToCurrency                        string  `json:"toCurrency"`
	RegularMarketChange               float64 `json:"regularMarketChange"`
	RegularMarketChangePercent        float64 `json:"regularMarketChangePercent"`
	RegularMarketTime                 int     `json:"regularMarketTime"`
	RegularMarketPrice                float64 `json:"regularMarketPrice"`
	RegularMarketDayHigh              float64 `json:"regularMarketDayHigh"`
	RegularMarketDayRange             string  `json:"regularMarketDayRange"`
	RegularMarketDayLow               float64 `json:"regularMarketDayLow"`
	RegularMarketVolume               int64   `json:"regularMarketVolume"`
	RegularMarketPreviousClose        float64 `json:"regularMarketPreviousClose"`
	FullExchangeName                  string  `json:"fullExchangeName"`
	RegularMarketOpen                 float64 `json:"regularMarketOpen"`
	AverageDailyVolume3Month          int64   `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day           int64   `json:"averageDailyVolume10Day"`
	StartDate                         int     `json:"startDate"`
	CoinImageURL                      string  `json:"coinImageUrl"`
	FiftyTwoWeekLowChange             float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent      float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekRange                 string  `json:"fiftyTwoWeekRange"`
	FiftyTwoWeekHighChange            float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent     float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow                   float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh                  float64 `json:"fiftyTwoWeekHigh"`
	FiftyDayAverage                   float64 `json:"fiftyDayAverage"`
	FiftyDayAverageChange             float64 `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent      float64 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage              float64 `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange        float64 `json:"twoHundredDayAverageChange"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	MarketCap                         int64   `json:"marketCap"`
	SourceInterval                    int     `json:"sourceInterval"`
	ExchangeDataDelayedBy             int     `json:"exchangeDataDelayedBy"`
	Tradeable                         bool    `json:"tradeable"`
	MarketState                       string  `json:"marketState"`
	Symbol                            string  `json:"symbol"`
}

// QuoteResponse is the top level json element from Yahoo Finance API
type QuoteResponse struct {
	Result []Result    `json:"result"`
	Error  interface{} `json:"error"`
}

// GetDetails returns the base summary of a ticker from Yahoo Finance API
func GetDetails(ticker string) *Result {
	res, err := http.Get(fmt.Sprintf(YahooURL, ticker))
	if err != nil {
		println(err)
		return nil
	}
	var out TLR
	d := json.NewDecoder(res.Body)
	d.Decode(&out)
	if err != nil || len(out.QuoteResponse.Result) == 0 {
		if err != nil {
			println(err)
		}
		return nil
	}
	return &out.QuoteResponse.Result[0]
}

// GetPrice retrieves the current trading price of a ticker from Yahoo
func GetPrice(ticker string) float64 {
	det := GetDetails(ticker)
	if det == nil {
		return -1
	}
	return det.RegularMarketPrice
}
