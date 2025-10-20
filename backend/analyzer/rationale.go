// Contiene la lógica de generación de explicaciones
package analyzer

import (
	"fmt"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

// generateRationale construye una explicación detallada del análisis realizado,
// incluyendo estadísticas clave y cambios significativos en las recomendaciones.
func (a *StockAnalyzer) generateRationale(stockHistory []models.Stock) string {
	if len(stockHistory) == 0 {
		return "Sin datos históricos disponibles para análisis"
	}

	latest := stockHistory[0]

	// Manejar brokerage vacío
	brokerage := latest.Brokerage
	if brokerage == "" {
		brokerage = "fuente no especificada"
	}

	// Manejar acción vacía
	action := latest.Action
	if action == "" {
		action = "acción no especificada"
	}

	// Manejar ratings vacíos
	ratingFrom := latest.RatingFrom
	if ratingFrom == "" {
		ratingFrom = "sin rating previo"
	}

	ratingTo := latest.RatingTo
	if ratingTo == "" {
		ratingTo = "rating no especificado"
	}

	// Manejar precios objetivo
	targetFrom, errFrom := a.formatter.ExtractNumber(latest.TargetFrom)
	targetTo, errTo := a.formatter.ExtractNumber(latest.TargetTo)

	targetText := ""
	if errFrom == nil && errTo == nil && targetFrom > 0 && targetTo > 0 {
		changePercent := ((targetTo - targetFrom) / targetFrom) * 100
		targetText = fmt.Sprintf("Precio objetivo actualizado de $%.2f a $%.2f (cambio de %.2f%%).",
			targetFrom, targetTo, changePercent)
	} else if errTo == nil && targetTo > 0 {
		targetText = fmt.Sprintf("Precio objetivo establecido en $%.2f.", targetTo)
	} else {
		targetText = "Precio objetivo no disponible."
	}

	return fmt.Sprintf(
		"Análisis basado en %d recomendaciones. "+
			"Última actualización por %s con acción '%s', cambiando calificación de '%s' a '%s'. "+
			"%s",
		len(stockHistory),
		brokerage,
		action,
		ratingFrom,
		ratingTo,
		targetText,
	)
}
