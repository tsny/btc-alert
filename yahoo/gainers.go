package yahoo

import (
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
	tickersQuery  = ".Fw(600)"
	allLinksQuery = "a"
	allRowsQuery  = "tr"
)

func getData(gainers bool) *goquery.Document {
	// Request the HTML page.
	url := losersURL
	if gainers {
		url = gainersURL
	}
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

// GetTableHeader returns a string array structured around the row returned by Yahoo Finance
func GetTableHeader() []string {
	return []string{"Symbol", "Name", "Price", "Change", "PercentChange",
		"Volume", "Avg Vol(3 Month)", "Market Cap", "PE Ratio"}
}

// GetTopMoversTickers returns the 25 top daily gainers tickers
func GetTopMoversTickers(gainers bool) []string {
	doc := getData(gainers)
	var tickers []string
	doc.Find(allLinksQuery).Each(func(i int, rows *goquery.Selection) {
		if rows.HasClass("Fw(600) C($linkColor)") && len(rows.Text()) <= 4 {
			tickers = append(tickers, rows.Text())
		}
	})
	return tickers
}

// GetTopMoversAsArray returns the top gainers of the day as a nested array
// Useful for putting into a table
func GetTopMoversAsArray(gainers bool) [][]string {
	doc := getData(gainers)
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
