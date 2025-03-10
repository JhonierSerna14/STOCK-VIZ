package api

import (
	"encoding/json"
	"net/http"

	"github.com/JhonierSerna14/STOCK-VIZ/database"
	"github.com/JhonierSerna14/STOCK-VIZ/service"
	"gorm.io/gorm"
)

type API struct {
	stockService *service.StockService
}

func NewAPI(db *gorm.DB) *API {
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

func (a *API) getAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := a.stockService.GetAllStocks()
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

func (a *API) deleteAllStocks(w http.ResponseWriter, r *http.Request) {
	if err := a.stockService.DeleteAllStocks(); err != nil {
		http.Error(w, "Error al borrar los stocks", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Todos los stocks han sido eliminados exitosamente",
	})
}

func (a *API) getRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := a.stockService.GetRecommendations()
	if err != nil {
		http.Error(w, "Error analyzing stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
