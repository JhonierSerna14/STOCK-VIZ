package calculator

import (
	"fmt"
	"math"
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

// CalculateScore calcula la puntuación total con análisis temporal histórico
func (c *ScoreCalculator) CalculateScore(stockHistory []models.Stock) float64 {
	if len(stockHistory) == 0 {
		return 0
	}

	latestStock := stockHistory[0]

	ratingScore := c.calculateRatingScore(stockHistory)
	targetScore := c.calculateTargetScore(latestStock)
	targetTrendScore := c.calculateTargetTrendScore(stockHistory)
	volatilityScore := c.calculateVolatilityScore(stockHistory)
	recencyScore := c.calculateRecencyScore(latestStock.Time)

	return ((ratingScore * c.weights.Rating) +
		(targetScore * c.weights.Target) +
		(targetTrendScore * c.weights.TargetTrend) +
		(volatilityScore * c.weights.Volatility) +
		(recencyScore * c.weights.Recency)) * 100
}

// calculateRatingScore evalúa la recomendación basada en ratings con análisis de tendencias históricas
func (c *ScoreCalculator) calculateRatingScore(history []models.Stock) float64 {
	if len(history) == 0 {
		return 0
	}

	latest := history[0]
	latestRating := strings.ToLower(latest.RatingTo)

	latestValue, ok := scoring.RatingMapping[latestRating]
	if !ok && latestRating != "" {
		// loguear rating desconocido para revisión futura
		logUnknownRating(latestRating)
		latestValue = scoring.RatingMapping["unknown"]
	}

	// Analizar tendencia histórica de ratings (últimos 6 meses)
	trendBonus := c.calculateRatingTrend(history)

	score := latestValue + trendBonus

	// Normalizar entre 0 y 1
	if score > 1.0 {
		score = 1.0
	} else if score < 0.0 {
		score = 0.0
	}

	return score
}

// calculateRatingTrend analiza la tendencia de mejora/deterioro en los ratings históricos
func (c *ScoreCalculator) calculateRatingTrend(history []models.Stock) float64 {
	if len(history) < 3 {
		return c.calculateImmediateChange(history)
	}

	// Analizar últimos 6 meses o hasta 10 registros
	maxRecords := 10
	if len(history) < maxRecords {
		maxRecords = len(history)
	}

	var ratingValues []float64
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	for i := 0; i < maxRecords && i < len(history); i++ {
		stock := history[i]
		if stock.Time.Before(sixMonthsAgo) {
			break
		}

		rating := strings.ToLower(stock.RatingTo)
		if value, exists := scoring.RatingMapping[rating]; exists {
			ratingValues = append(ratingValues, value)
		}
	}

	if len(ratingValues) < 3 {
		return c.calculateImmediateChange(history)
	}

	// Calcular tendencia usando regresión lineal simple
	// y = (n * Σ(xy) - Σx * Σy) / (n * Σ(x²) - (Σx)²)
	n := float64(len(ratingValues))
	var sumX, sumY, sumXY, sumX2 float64

	for i, rating := range ratingValues {
		x := float64(i)
		sumX += x
		sumY += rating
		sumXY += x * rating
		sumX2 += x * x
	}

	// Pendiente de la regresión lineal
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)

	// Convertir pendiente a bonus/penalización (-0.15 a +0.15)
	trendBonus := math.Max(-0.15, math.Min(0.15, slope*0.3))

	return trendBonus
}
func (c *ScoreCalculator) calculateImmediateChange(history []models.Stock) float64 {
	// Si no hay suficientes puntos, usar el cambio inmediato
	if len(history) >= 1 {
		latest := history[0]
		latestRating := strings.ToLower(latest.RatingTo)
		previousRating := strings.ToLower(latest.RatingFrom)

		latestValue, ok1 := scoring.RatingMapping[latestRating]
		if !ok1 {
			logUnknownRating(latestRating)
			latestValue = scoring.RatingMapping["unknown"]
		}

		previousValue, ok2 := scoring.RatingMapping[previousRating]
		if !ok2 {
			logUnknownRating(previousRating)
			previousValue = scoring.RatingMapping["unknown"]
		}

		// Bonus basado en el cambio inmediato
		immediateChange := latestValue - previousValue
		return math.Max(-0.15, math.Min(0.15, immediateChange))
	}
	return 0 // Si no hay datos suficientes, retornar 0
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

// calculateTargetTrendScore evalúa la evolución histórica de precios objetivo
func (c *ScoreCalculator) calculateTargetTrendScore(history []models.Stock) float64 {
	if len(history) < 2 {
		return 0.5 // Valor neutral si no hay suficiente historial
	}

	var targetChanges []float64
	sixMonthsAgo := time.Now().AddDate(0, -6, 0)

	for _, stock := range history {
		if stock.Time.Before(sixMonthsAgo) {
			break
		}

		targetFrom, errFrom := extractNumber(stock.TargetFrom)
		targetTo, errTo := extractNumber(stock.TargetTo)

		if errFrom == nil && errTo == nil && targetFrom > 0 && targetTo > 0 {
			change := (targetTo - targetFrom) / targetFrom
			targetChanges = append(targetChanges, change)
		}
	}

	if len(targetChanges) < 2 {
		return 0.5
	}

	// Calcular tendencia general de los cambios de target
	var avgChange float64
	for _, change := range targetChanges {
		avgChange += change
	}
	avgChange /= float64(len(targetChanges))

	// Calcular consistencia en la dirección de los cambios
	var positiveChanges int
	for _, change := range targetChanges {
		if change > 0.02 { // Cambios mayores al 2% se consideran significativos
			positiveChanges++
		}
	}

	// Proporción de cambios positivos significativos en relación con el total de cambios
	// Calcula el porcentaje de cambios positivos significativos en los precios objetivo.
	// Este valor se utiliza para medir la consistencia en la dirección de los cambios.
	// Un valor mayor a 0.5 indica que más de la mitad de los cambios son positivos.
	positiveRatio := float64(positiveChanges) / float64(len(targetChanges))

	// Combinar tendencia promedio y consistencia
	// Calcula un puntaje basado en la tendencia promedio de los cambios y la consistencia en la dirección de los mismos.
	// Escala el impacto de la tendencia promedio para que tenga un peso significativo en el puntaje final.
	// El valor base es 0.5, y se ajusta según el promedio de los cambios (avgChange).
	trendScore := 0.5 + (avgChange * 2.5)

	// Calcula un bonus o penalización basado en la consistencia de los cambios positivos.
	// Si más del 50% de los cambios son positivos (positiveRatio > 0.5), se otorga un bonus positivo.
	// Si menos del 50% de los cambios son positivos (positiveRatio < 0.5), se aplica una penalización.
	// El impacto del bonus está limitado a un rango de [-0.15, +0.15].
	consistencyBonus := (positiveRatio - 0.5) * 0.3

	// Puntaje final que combina la tendencia promedio y la consistencia.
	// Este valor representa una evaluación equilibrada de la evolución histórica de los precios objetivo.
	// Se suma el puntaje de la tendencia (trendScore) con el bonus de consistencia (consistencyBonus).
	finalScore := trendScore + consistencyBonus

	// Normalizar entre 0 y 1
	if finalScore > 1.0 {
		finalScore = 1.0
	} else if finalScore < 0.0 {
		finalScore = 0.0
	}

	return finalScore
}

// calculateVolatilityScore evalúa la consistencia histórica de las recomendaciones
func (c *ScoreCalculator) calculateVolatilityScore(history []models.Stock) float64 {
	if len(history) < 3 {
		return 0.5 // Valor neutral si no hay suficiente historial
	}

	var ratingValues []float64
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	// Recopilar valores de rating de los últimos 3 meses
	for _, stock := range history {
		if stock.Time.Before(threeMonthsAgo) {
			break
		}

		rating := strings.ToLower(stock.RatingTo)
		if value, exists := scoring.RatingMapping[rating]; exists {
			ratingValues = append(ratingValues, value)
		}
	}

	if len(ratingValues) < 3 {
		return 0.5
	}

	// Calcular volatilidad (desviación estándar)
	// La volatilidad mide la dispersión de los valores de los ratings en relación con su promedio.
	// Se calcula utilizando la desviación estándar, que es la raíz cuadrada de la varianza.
	// Fórmula de la varianza:
	// variance = Σ((xᵢ - mean)²) / n
	// Fórmula de la desviación estándar:
	// volatility = √variance

	// Calcular la media (mean) de los valores de los ratings
	var mean float64
	for _, value := range ratingValues {
		mean += value
	}
	mean /= float64(len(ratingValues)) // mean = Σ(xᵢ) / n

	// Calcular la varianza (variance)
	// La varianza mide el promedio de las diferencias al cuadrado entre cada valor y la media.
	var variance float64
	for _, value := range ratingValues {
		diff := value - mean    // Diferencia entre el valor actual y la media
		variance += diff * diff // Sumar el cuadrado de la diferencia
	}
	variance /= float64(len(ratingValues)) // variance = Σ((xᵢ - mean)²) / n

	// Calcular la desviación estándar (volatilidad)
	// La desviación estándar es la raíz cuadrada de la varianza.
	volatility := math.Sqrt(variance) // volatility = √variance

	// Convertir volatilidad a score: menor volatilidad = mayor score
	// Volatilidad típica está entre 0 y 0.3, normalizamos inversamente
	consistencyScore := 1.0 - math.Min(volatility/0.3, 1.0)

	// Agregar bonus por número de observaciones
	observationBonus := math.Min(float64(len(ratingValues))/10.0, 0.1)

	finalScore := consistencyScore + observationBonus

	// Normalizar entre 0 y 1
	if finalScore > 1.0 {
		finalScore = 1.0
	} else if finalScore < 0.0 {
		finalScore = 0.0
	}

	return finalScore
}

// calculateRecencyScore evalúa la actualidad de la información
func (c *ScoreCalculator) calculateRecencyScore(lastUpdate time.Time) float64 {
	daysAgo := time.Since(lastUpdate).Hours() / 24
	//función de decaimiento exponencial; el divisor 30.0 representa la "vida media" en días
	score := 1.0 / (1.0 + (daysAgo / 30.0))
	return score
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
