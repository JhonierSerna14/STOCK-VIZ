package models

import (
	"time"
)

// RecommendationFilter define los criterios para filtrar recomendaciones
type RecommendationFilter struct {
	Limit    int    // Número máximo de recomendaciones a retornar
	DateFrom string // Fecha desde (formato YYYY-MM-DD)
	DateTo   string // Fecha hasta (formato YYYY-MM-DD)
	Rating   string // Tipo de rating (buy, sell, hold, etc.)
	Ticker   string // Símbolo específico de la acción
}

// ParseDates convierte las cadenas de fechas a objetos time.Time
// Retorna fechas nil si las cadenas están vacías o mal formateadas
func (f *RecommendationFilter) ParseDates() (from, to *time.Time, err error) {
	if f.DateFrom != "" {
		parsedFrom, err := time.Parse("2006-01-02", f.DateFrom)
		if err != nil {
			return nil, nil, err
		}
		from = &parsedFrom
	}

	if f.DateTo != "" {
		parsedTo, err := time.Parse("2006-01-02", f.DateTo)
		if err != nil {
			return from, nil, err
		}
		// Ajustar al final del día
		parsedTo = parsedTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		to = &parsedTo
	}

	return from, to, nil
}
