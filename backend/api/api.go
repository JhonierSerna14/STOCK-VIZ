// Package api implementa los endpoints HTTP y la lógica de manejo de solicitudes.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/JhonierSerna14/STOCK-VIZ/database"
	"github.com/JhonierSerna14/STOCK-VIZ/service"
	"gorm.io/gorm"
)

// API encapsula los controladores HTTP y servicios necesarios para manejar las solicitudes.
type API struct {
	// stockService gestiona la lógica de negocio relacionada con las acciones
	stockService *service.StockService
}

// NewAPI crea y devuelve una nueva instancia de API configurada con la conexión a la base de datos.
// Inicializa el repositorio y servicio de stocks para ser utilizados por los controladores.
func NewAPI(db *gorm.DB) *API {
	stockRepo := database.NewStockRepository(db)
	return &API{
		stockService: service.NewStockService(stockRepo),
	}
}

// getStocks maneja solicitudes GET para obtener stocks con paginación.
// Acepta un parámetro de consulta 'next_page' para implementar paginación.
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

// getAllStocks maneja solicitudes GET para obtener todos los stocks sin paginación.
// Devuelve la lista completa de stocks en formato JSON.
func (a *API) getAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := a.stockService.GetAllStocks()
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

// deleteAllStocks maneja solicitudes DELETE para eliminar todos los stocks de la base de datos.
// Devuelve un mensaje de confirmación cuando la operación es exitosa.
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

// getRecommendations maneja solicitudes GET para obtener recomendaciones analizadas.
// Utiliza el servicio de stocks para generar recomendaciones basadas en datos existentes.
func (a *API) getRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := a.stockService.GetRecommendations()
	if err != nil {
		http.Error(w, "Error analyzing stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recommendations)
}
