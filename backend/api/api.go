package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/JhonierSerna14/STOCK-VIZ/database"
	"github.com/JhonierSerna14/STOCK-VIZ/service"
)

type API struct {
	stockService *service.StockService
}

func NewAPI(db *sql.DB) *API {
	stockRepo := database.NewStockRepository(db)
	return &API{
		stockService: service.NewStockService(stockRepo),
	}
}

func (a *API) getStocks(w http.ResponseWriter, r *http.Request) {
	nextPage := r.URL.Query().Get("next_page")

	stockResponse, err := a.stockService.GetStocks(nextPage)
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stockResponse)
}
