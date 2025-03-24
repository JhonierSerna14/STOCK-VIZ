// Package api implementa los endpoints HTTP y la lógica de manejo de solicitudes.
package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/JhonierSerna14/STOCK-VIZ/service"
)

// API encapsula los controladores HTTP y servicios necesarios para manejar las solicitudes.
type API struct {
	// stockService gestiona la lógica de negocio relacionada con las acciones
	stockService *service.StockService
}

// NewAPIWithService crea y devuelve una nueva instancia de API configurada con un servicio de stocks ya existente.
func NewAPIWithService(stockService *service.StockService) *API {
	return &API{
		stockService: stockService,
	}
}

// getStocks maneja solicitudes GET para obtener todos los stocks con paginación opcional.
// Acepta parámetros de consulta 'page' y 'limit' para implementar paginación. GET /api/stocks/all?page=2&limit=20
func (a *API) getStocks(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de paginación
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Convertir a enteros con valores predeterminados
	page := 1
	limit := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	query := r.URL.Query().Get("search")

	// Llamar al servicio con los parámetros de paginación
	stocks, total, err := a.stockService.GetAllStocksPaginated(page, limit, query)
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}

	// Preparar respuesta con metadatos de paginación
	response := map[string]interface{}{
		"items": stocks,
		"pagination": map[string]interface{}{
			"current_page": page,
			"per_page":     limit,
			"total_items":  total,
			"total_pages":  (total + int64(limit) - 1) / int64(limit),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
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
	err := json.NewEncoder(w).Encode(map[string]string{
		"mensaje": "Todos los stocks han sido eliminados exitosamente",
	})
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// getRecommendations maneja solicitudes GET para obtener recomendaciones analizadas.
// Acepta parámetros de consulta para filtrado y paginación.
func (a *API) getRecommendations(w http.ResponseWriter, r *http.Request) {
	// Obtener parámetros de filtrado
	query := r.URL.Query()

	// Parámetro de límite (por defecto 6)
	limitStr := query.Get("limit")
	limit := 6
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Fechas
	dateFrom := query.Get("date_from")
	dateTo := query.Get("date_to")

	// Tipo de rating
	rating := query.Get("rating")

	// Ticker específico
	ticker := query.Get("ticker")

	// Crear filtro
	filter := models.RecommendationFilter{
		Limit:    limit,
		DateFrom: dateFrom,
		DateTo:   dateTo,
		Rating:   rating,
		Ticker:   ticker,
	}

	recommendations, err := a.stockService.GetRecommendations(filter)
	if err != nil {
		http.Error(w, "Error analyzing stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(recommendations)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
