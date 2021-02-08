package yahoo

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Summary struct {
	Symbol        string
	Name          string
	Price         float64
	Change        float64
	PercentChange float64
	Volume        float64
}

const (
	url = "https://finance.yahoo.com/gainers/"
	// allLinksQuery = "a .Fw(600) .C($linkColor)"
	allLinksQuery = "a"
	allRowsQuery  = "tr"
)

func getData() *goquery.Document {
	// Request the HTML page.
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

func GetTopGainersAsArray() [][]string {
	doc := getData()

	var arr [][]string
	// Symbol, Name, Price, Change, % Change, Volume
	println(strings.Join(GetTableHeader(), " | "))
	doc.Find(allRowsQuery).Each(func(i int, rows *goquery.Selection) {
		cells := rows.Find("td")
		str := ""
		inner := []string{}
		cells.Each(func(j int, cells *goquery.Selection) {
			inner = append(inner, cells.Text())
			str += " | " + cells.Text()
		})
		arr = append(arr, inner)
		println(str)
	})
	return arr
}
