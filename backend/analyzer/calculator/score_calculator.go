package calculator

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/formatter"
	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/scoring"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

// ScoreCalculator maneja los cálculos de puntuación
type ScoreCalculator struct {
	weights   scoring.FactorWeights
	formatter *formatter.NumberFormatter
}

// NewScoreCalculator crea una nueva instancia de ScoreCalculator
func NewScoreCalculator(weights scoring.FactorWeights) *ScoreCalculator {
	return &ScoreCalculator{
		weights:   weights,
		formatter: &formatter.NumberFormatter{},
	}
}

// CalculateScore calcula la puntuación total
func (c *ScoreCalculator) CalculateScore(stockHistory []models.Stock) float64 {
	if len(stockHistory) == 0 {
		return 0
	}

	latestStock := stockHistory[0]

	ratingScore := c.calculateRatingScore(stockHistory)
	targetScore := c.calculateTargetScore(latestStock)
	brokerScore := c.calculateBrokerConsensus(stockHistory)
	recencyScore := c.calculateRecencyScore(latestStock.Time)

	return (ratingScore*c.weights.Rating +
		targetScore*c.weights.Target +
		brokerScore*c.weights.Broker +
		recencyScore*c.weights.Recency) * 100
}

// calculateRatingScore evalúa la recomendación basada en ratings
func (c *ScoreCalculator) calculateRatingScore(history []models.Stock) float64 {
	latest := history[0]
	latestRating := strings.ToLower(latest.RatingTo)
	previousRating := strings.ToLower(latest.RatingFrom)

	latestValue, ok1 := scoring.RatingMapping[latestRating]
	if !ok1 {
		logUnknownRating(latestRating)
		latestValue = 0.5
	}
	previousValue, ok2 := scoring.RatingMapping[previousRating]
	if !ok2 {
		logUnknownRating(previousRating)
		previousValue = 0.5
	}

	score := latestValue
	if latestValue > previousValue {
		score += 0.1
	}
	if score > 1.0 {
		score = 1.0
	} else if score < 0.0 {
		score = 0.0
	}
	return score
}

// calculateTargetScore evalúa el cambio en el precio objetivo
func (c *ScoreCalculator) calculateTargetScore(stock models.Stock) float64 {
	targetFrom, errFrom := extractNumber(stock.TargetFrom)
	targetTo, errTo := extractNumber(stock.TargetTo)

	if errFrom != nil || errTo != nil {
		return 0
	}
	if targetFrom <= 0 || targetTo <= 0 {
		return 0
	}

	changePercent := ((targetTo - targetFrom) / targetFrom) * 100

	// Normalizamos el cambio porcentual a un score entre 0 y 1
	// Un cambio de -20% o menor es 0, un cambio de +20% o mayor es 1
	score := changePercent / 20
	if score > 1.0 {
		score = 1.0
	} else if score < 0.0 {
		score = 0.0
	}
	return score
}

// calculateBrokerConsensus evalúa el consenso de los brokers
func (c *ScoreCalculator) calculateBrokerConsensus(history []models.Stock) float64 {
	brokerData := make(map[string]struct {
		sentimentSum float64
		count        int
	})

	for _, stock := range history {
		action := normalizeAction(stock.Action)
		var sentiment float64
		switch action {
		case "upgraded", "target raised", "initiated":
			sentiment = 1
		case "downgraded", "target lowered":
			sentiment = -1
		case "reiterated", "maintained":
			sentiment = 0
		}
		entry := brokerData[stock.Brokerage]
		entry.sentimentSum += sentiment
		entry.count++
		brokerData[stock.Brokerage] = entry
	}

	var totalWeightedSentiment float64
	var totalCount int
	for _, entry := range brokerData {
		// Calcula el promedio de sentimiento; Pesa por el número de recomendaciones
		totalWeightedSentiment += (entry.sentimentSum / float64(entry.count)) * float64(entry.count)
		totalCount += entry.count
	}

	if totalCount == 0 {
		return 0
	}
	//El score siempre esté normalizado entre 0 y 1
	normalizedScore := (totalWeightedSentiment + float64(totalCount)) / (2 * float64(totalCount))
	return normalizedScore
}

// calculateRecencyScore evalúa la actualidad de la información
func (c *ScoreCalculator) calculateRecencyScore(lastUpdate time.Time) float64 {
	daysAgo := time.Since(lastUpdate).Hours() / 24
	//función de decaimiento exponencial; el divisor 30.0 representa la "vida media" en días
	score := 1.0 / (1.0 + (daysAgo / 30.0))
	return score
}

// Funciones auxiliares
func normalizeAction(action string) string {
	return strings.ToLower(strings.TrimSpace(action))
}

func extractNumber(value string) (float64, error) {
	clean := strings.ReplaceAll(value, ",", "")
	clean = strings.TrimSpace(strings.Replace(clean, "$", "", -1))
	return strconv.ParseFloat(clean, 64)
}

// logUnknownRating registra un rating desconocido en el archivo de logs
func logUnknownRating(rating string) {
	f, err := os.OpenFile("logs/unknown_ratings.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			return
		}
	}()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("%s - %s\n", timestamp, rating)

	if _, err := f.WriteString(logEntry); err != nil {
		return
	}
}
