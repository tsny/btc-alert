package stocks

import (
	"context"
	"strings"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
	"github.com/sirupsen/logrus"
)

var PC *polygon.Client

// GetPrice retrieves the current trading price of a ticker from Yahoo
func GetPrice(ticker string) float64 {
	ticker = strings.ToUpper(ticker)
	params := &models.GetLastTradeParams{
		Ticker: ticker,
	}

	res, err := PC.GetLastTrade(context.Background(), params)
	if err != nil {
		logrus.Errorf(err.Error())
		return -1
	}

	return res.Results.Price
}
