package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *API) RegisterRoutes(r *mux.Router) {

	r.HandleFunc("/api/stocks", a.getStocks).Methods(http.MethodGet)
	r.HandleFunc("/api/stocks/all", a.getAllStocks).Methods(http.MethodGet)
	r.HandleFunc("/api/stocks", a.deleteAllStocks).Methods(http.MethodDelete)
	r.HandleFunc("/api/recommendations", a.getRecommendations).Methods(http.MethodGet)
	r.HandleFunc("/api/stocks/migrate", a.migrateAllStocks).Methods(http.MethodGet)
}
