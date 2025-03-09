package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *API) RegisterRoutes(r *mux.Router) {

	r.HandleFunc("/api/stocks", a.getStocks).Methods(http.MethodGet)
}
