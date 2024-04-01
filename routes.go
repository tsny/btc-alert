package main

import (
	"encoding/json"
	"net/http"

	"btc-alert/eps"

	"github.com/gorilla/mux"
)

// TODO: we could have the GET symbol endpoints return references to
// other endpoints like prices/graph?
func initRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/symbol/{symbol}/prices", getRecentPrices).Methods("GET")
	r.HandleFunc("/all", getAll).Methods("GET")
	return r
}

func getAll(w http.ResponseWriter, r *http.Request) {
	var arr []*eps.Publisher
	for _, v := range lookupService.GetAllTracked() {
		arr = append(arr, v.Publisher)
	}
	sendJSON(w, arr)
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
