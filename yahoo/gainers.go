package yahoo

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

/*
	This section of the package deals with querying Yahoo Finance's
	Chart table for various criteria, mainly top gainers
	We scrape the URL below and grab HTML elements from it
	because I cannot find the URL for the API if it still exists
*/

type Summary struct {
	Symbol        string
	Name          string
	Price         float64
	Change        float64
	PercentChange float64
	Volume        float64
}

const (
	gainersURL    = "https://finance.yahoo.com/gainers/"
	losersURL     = "https://finance.yahoo.com/losers/"
	baseURL       = "https://www.marketwatch.com/investing/stock/%s"
	tickersQuery  = ".Fw(600)"
	allLinksQuery = "a"
	allRowsQuery  = "tr"
	// summaryQuery  = "p.businessSummary"
	summaryQuery = ".businessSummary.Mt(10px).Ov(h).Tov(e)"
)

func getData(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return doc
}

// GetSummary returns the business summary of a ticker from yahoo finance
func GetSummary(ticker string) string {
	url := fmt.Sprintf(baseURL, ticker)
	doc := getData(url)
	return doc.Find(".description__text").Text()
}

// GetTableHeader returns a string array structured around the row returned by Yahoo Finance
func GetTableHeader() []string {
	return []string{"Symbol", "Name", "Price", "Change", "PercentChange",
		"Volume", "Avg Vol(3 Month)", "Market Cap", "PE Ratio"}
}

// GetTopMoversTickers returns the 25 top daily gainers tickers
func GetTopMoversTickers(gainers bool) []string {
	url := losersURL
	if gainers {
		url = gainersURL
	}
	doc := getData(url)
	var tickers []string
	doc.Find(allLinksQuery).Each(func(i int, rows *goquery.Selection) {
		if rows.HasClass("Fw(600) C($linkColor)") && len(rows.Text()) <= 4 {
			tickers = append(tickers, rows.Text())
		}
	})
	return tickers
}

func GetGainers() []string {
	doc := getData(gainersURL)
	var arr []string
	doc.Find("[data-test='quoteLink']").Each(func(i int, ele *goquery.Selection) {
		arr = append(arr, ele.Text())
	})
	if len(arr) > 10 {
		arr = arr[:9]
	}
	return arr
}

// GetTopMoversAsArray returns the top gainers of the day as a nested array
// Useful for putting into a table
func GetTopMoversAsArray(gainers bool) [][]string {
	url := losersURL
	if gainers {
		url = gainersURL
	}
	doc := getData(url)
	var arr [][]string
	doc.Find(allRowsQuery).Each(func(i int, rows *goquery.Selection) {
		cells := rows.Find("td")
		str := ""
		inner := []string{}
		cells.Each(func(j int, cells *goquery.Selection) {
			inner = append(inner, cells.Text())
			str += " | " + cells.Text()
		})
		arr = append(arr, inner)
	})
	return arr
}
