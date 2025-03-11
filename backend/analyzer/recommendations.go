// Contiene la lógica de recomendaciones
package analyzer

import (
	"sort"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

func (a *StockAnalyzer) GetTopRecommendations(limit int) ([]models.StockRecommendation, error) {
	stocks, err := a.repository.GetAllStocks()
	if err != nil {
		return nil, err
	}

	stockMap := a.groupStocksByTicker(stocks)
	recommendations := a.generateRecommendations(stockMap)

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if limit > 0 && limit < len(recommendations) {
		recommendations = recommendations[:limit]
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
