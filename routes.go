package main

import (
	"encoding/json"
	"net/http"

	"btc-alert/eps"
	"btc-alert/yahoo"

	"github.com/gorilla/mux"
)

// TODO: we could have the GET symbol endpoints return references to
// other endpoints like prices/graph?
func initRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/symbol/{symbol}", getSymbol).Methods("GET")
	r.HandleFunc("/symbol/{symbol}/prices", getRecentPrices).Methods("GET")
	r.HandleFunc("/symbol/{symbol}/details", getStockDetails).Methods("GET")
	r.HandleFunc("/gainers", getLosers).Methods("GET")
	r.HandleFunc("/losers", getGainers).Methods("GET")
	r.HandleFunc("/all", getAll).Methods("GET")
	// r.HandleFunc("/watchlist", postRefreshWatchlist).Methods("POST")
	return r
}

func getAll(w http.ResponseWriter, r *http.Request) {
	var arr []*eps.Publisher
	for _, v := range lookupService.GetAllTracked() {
		arr = append(arr, v.Publisher)
	}
	sendJSON(w, arr)
}

func getGainers(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, yahoo.GetTopMoversAsArray(true))
}

func getLosers(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, yahoo.GetTopMoversAsArray(false))
}

func getSymbol(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if sec := lookupService.FindSecurityByNameOrTicker(symbol); sec != nil {
		sendJSON(w, sec.Publisher)
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
	if sec := lookupService.FindSecurityByNameOrTicker(symbol); sec != nil {
		sendJSON(w, sec.Queue.GetAllPrices())
	}
}

func sendJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
