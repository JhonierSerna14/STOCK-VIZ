// Contiene la lógica de recomendaciones
package analyzer

import (
	"sort"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

func (a *StockAnalyzer) GetFilteredRecommendations(filter models.RecommendationFilter) ([]models.StockRecommendation, error) {
	// Obtener stocks filtrados directamente de la base de datos
	stocks, err := a.repository.GetFilteredStocks(filter)
	if err != nil {
		return nil, err
	}

	// Agrupar stocks por ticker
	stockMap := a.groupStocksByTicker(stocks)
	// Generar recomendaciones
	recommendations := a.generateRecommendations(stockMap)

	// Ordenar por puntuación
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Aplicar límite si es necesario
	if filter.Limit > 0 && filter.Limit < len(recommendations) {
		recommendations = recommendations[:filter.Limit]
	}

	return recommendations, nil
}

func (a *StockAnalyzer) generateRecommendations(stockMap map[string][]models.Stock) []models.StockRecommendation {
	var recommendations []models.StockRecommendation

	// Asegura que las recomendaciones más recientes estén primero
	for ticker, stockHistory := range stockMap {
		sort.Slice(stockHistory, func(i, j int) bool {
			return stockHistory[i].Time.After(stockHistory[j].Time)
		})

		latestStock := stockHistory[0]
		score := a.calculator.CalculateScore(stockHistory)

		recommendation := models.StockRecommendation{
			Ticker:            ticker,
			Company:           latestStock.Company,
			Score:             score,
			LatestRating:      latestStock.RatingTo,
			LatestTarget:      latestStock.TargetTo,
			LastUpdated:       latestStock.Time,
			AnalysisRationale: a.generateRationale(stockHistory),
		}

		recommendations = append(recommendations, recommendation)
	}

	return recommendations
}

func (a *StockAnalyzer) groupStocksByTicker(stocks []models.Stock) map[string][]models.Stock {
	stockMap := make(map[string][]models.Stock)
	for _, stock := range stocks {
		stockMap[stock.Ticker] = append(stockMap[stock.Ticker], stock)
	}
	return stockMap
}
