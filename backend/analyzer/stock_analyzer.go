package analyzer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

type StockAnalyzer struct {
	repository StockRepository
}

type StockRepository interface {
	GetAllStocks() ([]models.Stock, error)
}

type StockRecommendation struct {
	Ticker            string    `json:"ticker"`
	Company           string    `json:"company"`
	Score             float64   `json:"score"`
	LatestRating      string    `json:"latest_rating"`
	LatestTarget      string    `json:"latest_target"`
	LastUpdated       time.Time `json:"last_updated"`
	AnalysisRationale string    `json:"analysis_rationale"`
}

func NewStockAnalyzer(repo StockRepository) *StockAnalyzer {
	return &StockAnalyzer{repository: repo}
}

func (a *StockAnalyzer) GetTopRecommendations(limit int) ([]StockRecommendation, error) {
	stocks, err := a.repository.GetAllStocks()
	if err != nil {
		return nil, err
	}

	// Agrupar stocks por ticker para análisis
	stockMap := make(map[string][]models.Stock)
	for _, stock := range stocks {
		stockMap[stock.Ticker] = append(stockMap[stock.Ticker], stock)
	}

	var recommendations []StockRecommendation
	for ticker, stockHistory := range stockMap {
		// Ordenar el historial por tiempo (datos más recientes primero)
		sort.Slice(stockHistory, func(i, j int) bool {
			return stockHistory[i].Time.After(stockHistory[j].Time)
		})

		latestStock := stockHistory[0]
		score := a.calculateScore(stockHistory)

		recommendation := StockRecommendation{
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

	// Ordenar recomendaciones por puntuación (mayor a menor)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if limit > 0 && limit < len(recommendations) {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

func (a *StockAnalyzer) calculateScore(stockHistory []models.Stock) float64 {
	if len(stockHistory) == 0 {
		return 0
	}

	latestStock := stockHistory[0]

	// Factor 1: Análisis de calificaciones (25%)
	ratingScore := a.calculateRatingScore(stockHistory)

	// Factor 2: Análisis de precio objetivo (25%)
	targetScore := a.calculateTargetScore(latestStock)

	// Factor 3: Consenso de brokers (30%)
	brokerScore := a.calculateBrokerConsensus(stockHistory)

	// Factor 4: Recencia de las actualizaciones (20%)
	recencyScore := a.calculateRecencyScore(latestStock.Time)

	// Sumar los factores según sus pesos y normalizar a escala 0-100
	score := (ratingScore*0.25 + targetScore*0.25 + brokerScore*0.30 + recencyScore*0.20) * 100
	return score
}

func (a *StockAnalyzer) calculateRatingScore(history []models.Stock) float64 {
	// Mapa de ratings extendido para contemplar variaciones comunes
	ratingMapping := map[string]float64{
		"strong buy":          1.0,
		"buy":                 0.8,
		"outperform":          0.75,
		"overweight":          0.70,
		"equal weight":        0.5,
		"neutral":             0.5,
		"in-line":             0.5,
		"hold":                0.4,
		"market perform":      0.45,
		"sector outperform":   0.75,
		"sector perform":      0.5,
		"sector underperform": 0.3,
		"underweight":         0.3,
		"underperform":        0.3,
		"sell":                0.0,
	}

	latest := history[0]
	latestRating := strings.ToLower(latest.RatingTo)
	previousRating := strings.ToLower(latest.RatingFrom)

	// Obtener los valores numéricos; si no existe la clave, se asume 0
	latestValue, ok1 := ratingMapping[latestRating]
	if !ok1 {
		latestValue = 0.5
	}
	previousValue, ok2 := ratingMapping[previousRating]
	if !ok2 {
		previousValue = 0.5
	}

	score := latestValue

	// Bonus si hay mejora en la calificación
	if latestValue > previousValue {
		score += 0.1
	}

	// Asegurarse que el score esté en rango [0,1]
	if score > 1.0 {
		score = 1.0
	} else if score < 0.0 {
		score = 0.0
	}

	return score
}

func (a *StockAnalyzer) calculateTargetScore(stock models.Stock) float64 {
	targetFrom := a.extractNumber(stock.TargetFrom)
	targetTo := a.extractNumber(stock.TargetTo)

	if targetFrom <= 0 || targetTo <= 0 {
		return 0
	}

	// Calcular el cambio porcentual
	changePercent := ((targetTo - targetFrom) / targetFrom) * 100

	// Normalizar el score basado en el cambio porcentual:
	// Se asume que un cambio de +20% o más es muy positivo (score = 1.0)
	// y un cambio de -20% o menos es muy negativo (score = 0.0)
	score := (changePercent + 20) / 40
	if score > 1.0 {
		score = 1.0
	} else if score < 0.0 {
		score = 0.0
	}

	return score
}

func (a *StockAnalyzer) calculateBrokerConsensus(history []models.Stock) float64 {
	brokerSentiment := make(map[string]int)

	for _, stock := range history {
		action := normalizeAction(stock.Action)
		sentiment := 0
		switch action {
		case "upgraded", "target raised", "initiated":
			sentiment = 1
		case "downgraded", "target lowered":
			sentiment = -1
		case "reiterated", "maintained":
			sentiment = 0
		}
		brokerSentiment[stock.Brokerage] = sentiment
	}

	// Calcular el consenso: se promedian los sentimientos de cada broker
	var totalSentiment float64
	for _, sentiment := range brokerSentiment {
		totalSentiment += float64(sentiment)
	}

	brokerCount := float64(len(brokerSentiment))
	if brokerCount == 0 {
		return 0
	}

	// Normalizar el consenso a un rango de 0 a 1
	normalizedScore := (totalSentiment + brokerCount) / (2 * brokerCount)
	return normalizedScore
}

func (a *StockAnalyzer) calculateRecencyScore(lastUpdate time.Time) float64 {
	daysAgo := time.Since(lastUpdate).Hours() / 24
	// Se usa una función de decaimiento exponencial: actualizaciones recientes tienen score cercano a 1
	score := 1.0 / (1.0 + (daysAgo / 30.0))
	return score
}

func (a *StockAnalyzer) generateRationale(stockHistory []models.Stock) string {
	latest := stockHistory[0]
	brokerCount := len(stockHistory)

	targetFrom := a.extractNumber(latest.TargetFrom)
	targetTo := a.extractNumber(latest.TargetTo)
	targetChange := ""
	if targetFrom > 0 && targetTo > 0 {
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

func (a *StockAnalyzer) extractNumber(value string) float64 {
	// Remover el símbolo "$" y espacios
	clean := strings.TrimSpace(strings.Replace(value, "$", "", -1))
	num, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	return num
}

// normalizeAction estandariza la acción removiendo sufijos innecesarios y convirtiendo a minúsculas.
func normalizeAction(action string) string {
	norm := strings.ToLower(strings.TrimSpace(action))
	// Remover el sufijo " by" si existe
	norm = strings.TrimSuffix(norm, " by")
	return norm
}
