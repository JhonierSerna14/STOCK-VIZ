// StockRecommendation representa una recomendación completa para una acción específica,
// incluyendo métricas clave y justificación del análisis.

package models

import "time"

type StockRecommendation struct {
	Ticker            string    `json:"ticker"`             // Símbolo de la acción
	Company           string    `json:"company"`            // Nombre de la empresa
	Score             float64   `json:"score"`              // Puntuación calculada (0-100)
	LatestRating      string    `json:"latest_rating"`      // Última calificación asignada
	LatestTarget      string    `json:"latest_target"`      // Último precio objetivo
	LastUpdated       time.Time `json:"last_updated"`       // Fecha de última actualización
	AnalysisRationale string    `json:"analysis_rationale"` // Explicación del análisis
}
