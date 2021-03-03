package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsny/btc-alert/eps"
	"github.com/tsny/btc-alert/yahoo"
)

func initRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/symbol/{symbol}", getSymbol).Methods("GET")
	r.HandleFunc("/topMovers", getTopMovers).Methods("GET")
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
	if pub, ok := PublisherMap[symbol]; ok {
		sendJSON(w, pub)
	} else {
		sendJSON(w, yahoo.GetDetails(symbol))
	}
}

func getWatchlist(w http.ResponseWriter, r *http.Request) {
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
