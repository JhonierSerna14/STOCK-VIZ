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
	location, err := time.LoadLocation("Local")
	if err != nil {
		// Si no se puede cargar la zona horaria local, usamos UTC
		location = time.UTC
	}

	if f.DateFrom != "" {
		// Parse the date in the local timezone to avoid UTC conversion issues
		parsedFrom, err := time.ParseInLocation("2006-01-02", f.DateFrom, location)
		if err != nil {
			return nil, nil, err
		}
		from = &parsedFrom
	}

	if f.DateTo != "" {
		// Parse the date in the local timezone to avoid UTC conversion issues
		parsedTo, err := time.ParseInLocation("2006-01-02", f.DateTo, location)
		if err != nil {
			return from, nil, err
		}
		// Ajustar al final del día en la zona horaria local
		parsedTo = parsedTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		to = &parsedTo
	}
	return from, to, nil
}
