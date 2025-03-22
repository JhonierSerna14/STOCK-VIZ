// Contiene la lógica de recomendaciones
package analyzer

import (
	"sort"
	"strings"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

func (a *StockAnalyzer) GetFilteredRecommendations(filter models.RecommendationFilter) ([]models.StockRecommendation, error) {
	stocks, err := a.repository.GetAllStocks()
	if err != nil {
		return nil, err
	}

	// Aplicar filtro por ticker si está especificado
	if filter.Ticker != "" {
		var filteredStocks []models.Stock
		for _, stock := range stocks {
			if stock.Ticker == filter.Ticker {
				filteredStocks = append(filteredStocks, stock)
			}
		}
		stocks = filteredStocks
	}

	// Aplicar filtro por fechas si están especificadas
	dateFrom, dateTo, err := filter.ParseDates()
	if err != nil {
		return nil, err
	}

	if dateFrom != nil || dateTo != nil {
		var filteredStocks []models.Stock
		for _, stock := range stocks {
			include := true

			if dateFrom != nil && stock.Time.Before(*dateFrom) {
				include = false
			}

			if dateTo != nil && stock.Time.After(*dateTo) {
				include = false
			}

			if include {
				filteredStocks = append(filteredStocks, stock)
			}
		}
		stocks = filteredStocks
	}

	stockMap := a.groupStocksByTicker(stocks)
	recommendations := a.generateRecommendations(stockMap)

	// Aplicar filtro por rating si está especificado
	if filter.Rating != "" {
		var filteredRecommendations []models.StockRecommendation
		for _, rec := range recommendations {
			if strings.EqualFold(rec.LatestRating, filter.Rating) {
				filteredRecommendations = append(filteredRecommendations, rec)
			}
		}
		recommendations = filteredRecommendations
	}

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
