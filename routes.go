package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tsny/btc-alert/yahoo"
)

func initRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/symbol/{symbol}", getSymbol).Methods("GET")
	r.HandleFunc("/topMovers", getTopMovers).Methods("GET")
	return r
}

func getTopMovers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(yahoo.GetTopMoversAsArray(true))
}

func getSymbol(w http.ResponseWriter, r *http.Request) {
	symbol := mux.Vars(r)["symbol"]
	if pub, ok := PublisherMap[symbol]; ok {
		json.NewEncoder(w).Encode(pub)
	} else {
		json.NewEncoder(w).Encode(yahoo.GetPrice(symbol))
	}
}
