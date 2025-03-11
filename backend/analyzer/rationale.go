// Contiene la lógica de generación de explicaciones
package analyzer

import (
	"fmt"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

// generateRationale construye una explicación detallada del análisis realizado,
// incluyendo estadísticas clave y cambios significativos en las recomendaciones.
func (a *StockAnalyzer) generateRationale(stockHistory []models.Stock) string {
	latest := stockHistory[0]
	brokerCount := len(stockHistory)

	targetFrom, errFrom := a.formatter.ExtractNumber(latest.TargetFrom)
	targetTo, errTo := a.formatter.ExtractNumber(latest.TargetTo)
	targetChange := ""
	if errFrom == nil && errTo == nil && targetFrom > 0 && targetTo > 0 {
		changePercent := ((targetTo - targetFrom) / targetFrom) * 100
		targetChange = fmt.Sprintf(" (cambio de %.2f%%)", changePercent)
	}

	return fmt.Sprintf(
		"Análisis basado en %d recomendaciones de brokers. "+
			"Última actualización por %s con acción '%s', cambiando calificación de '%s' a '%s'. "+
			"Precio objetivo actualizado de $%.2f a $%.2f%s",
		brokerCount,
		latest.Brokerage,
		latest.Action,
		latest.RatingFrom,
		latest.RatingTo,
		targetFrom,
		targetTo,
		targetChange,
	)
}
