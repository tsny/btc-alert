package main

import (
	"encoding/json"
	"net/http"

	"btc-alert/eps"
	"btc-alert/priceTracking"
	"btc-alert/yahoo"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/wcharczuk/go-chart"
)

// TODO: we could have the GET symbol endpoints return references to
// other endpoints like prices/graph?
func initRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/symbol/{symbol}", getSymbol).Methods("GET")
	r.HandleFunc("/symbol/{symbol}/graph", getGraph).Methods("GET")
	r.HandleFunc("/symbol/{symbol}/prices", getRecentPrices).Methods("GET")
	r.HandleFunc("/symbol/{symbol}/details", getStockDetails).Methods("GET")
	r.HandleFunc("/movers", getTopMovers).Methods("GET")
	r.HandleFunc("/watchlist", getWatchlist).Methods("GET")
	r.HandleFunc("/watchlist", postRefreshWatchlist).Methods("POST")
	r.HandleFunc("/crypto", getCryptoWatchlist).Methods("GET")
	return r
}

func getTopMovers(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, yahoo.GetTopMoversAsArray(true))
}

func getSymbol(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if sec, pub := lookupService.FindSecurityByNameOrTicker(symbol); sec != nil {
		sendJSON(w, pub)
	} else {
		sendJSON(w, yahoo.GetDetails(symbol))
	}
}

func getStockDetails(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	sendJSON(w, yahoo.GetDetails(symbol))
}

func getRecentPrices(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	queue := queueService.FindByTicker(symbol)
	if queue == nil {
		logrus.Warnf("couldn't find ticker %s", symbol)
		return
	}
	sendJSON(w, queue.GetAllPrices())
}

func getGraph(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	queue := queueService.FindByTicker(symbol)
	if queue == nil {
		logrus.Warnf("couldn't find ticker %s", symbol)
		return
	}
	graph := priceTracking.QueueToGraph(*queue)
	w.Header().Set("Content-Type", "image/png")
	graph.Render(chart.PNG, w)
}

func getWatchlist(w http.ResponseWriter, r *http.Request) {
	if len(watchlist) == 0 {
		refreshWatchlist()
	}
	var publishers []*eps.Publisher
	for _, p := range watchlist {
		publishers = append(publishers, p)
	}
	sendJSON(w, publishers)
}

func getCryptoWatchlist(w http.ResponseWriter, r *http.Request) {
	var cr []*eps.Publisher
	for _, p := range PublisherMap {
		if p.UseMarketHours {
			continue
		}
		cr = append(cr, p)
	}
	sendJSON(w, cr)
}

func postRefreshWatchlist(w http.ResponseWriter, r *http.Request) {
	refreshWatchlist()
	getWatchlist(w, r)
}

func sendJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
