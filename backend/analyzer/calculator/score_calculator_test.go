package calculator

import (
	"os"
	"testing"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/scoring"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

func TestNewScoreCalculator(t *testing.T) {
	// Test creating a new calculator
	weights := scoring.DefaultWeights
	calculator := NewScoreCalculator(weights)

	if calculator == nil {
		t.Fatal("Expected calculator, got nil")
	}
	if calculator.formatter == nil {
		t.Error("Expected formatter to be initialized")
	}
}

func TestScoreCalculatorCalculateScoreEmptyHistory(t *testing.T) {
	// Test with empty stock history
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	score := calculator.CalculateScore([]models.Stock{})
	if score != 0.0 {
		t.Errorf("Expected score 0.0 for empty history, got %f", score)
	}
}

func TestScoreCalculatorCalculateScoreSingleStock(t *testing.T) {
	// Test with single stock
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	stocks := []models.Stock{
		{
			Ticker:   "AAPL",
			Company:  "Apple Inc",
			RatingTo: "buy",
			TargetTo: "175.00",
			Time:     time.Now(),
		},
	}

	score := calculator.CalculateScore(stocks)
	if score <= 0.0 {
		t.Errorf("Expected score > 0.0, got %f", score)
	}
	if score > 100.0 {
		t.Errorf("Expected score <= 100.0, got %f", score)
	}
}

func TestScoreCalculatorCalculateScoreMultipleStocks(t *testing.T) {
	// Test with multiple stocks showing upgrade trend
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	baseTime := time.Now()
	stocks := []models.Stock{
		{
			Ticker:     "AAPL",
			Company:    "Apple Inc",
			RatingFrom: "hold",
			RatingTo:   "buy",
			TargetFrom: "150.00",
			TargetTo:   "175.00",
			Time:       baseTime,
		},
		{
			Ticker:     "AAPL",
			Company:    "Apple Inc",
			RatingFrom: "sell",
			RatingTo:   "hold",
			TargetFrom: "130.00",
			TargetTo:   "150.00",
			Time:       baseTime.Add(-24 * time.Hour),
		},
	}

	score := calculator.CalculateScore(stocks)
	if score <= 0.0 {
		t.Errorf("Expected score > 0.0, got %f", score)
	}
	if score > 100.0 {
		t.Errorf("Expected score <= 100.0, got %f", score)
	}
}

func TestScoreCalculatorCalculateScoreCustomWeights(t *testing.T) {
	// Test with custom weights
	customWeights := scoring.FactorWeights{
		Rating:      1.0, // Only rating matters
		Target:      0.0,
		TargetTrend: 0.0,
		Volatility:  0.0,
		Recency:     0.0,
	}

	calculator := NewScoreCalculator(customWeights)

	stocks := []models.Stock{
		{
			Ticker:   "AAPL",
			RatingTo: "strong buy", // Should give high score
			Time:     time.Now(),
		},
	}

	score := calculator.CalculateScore(stocks)
	if score < 80.0 {
		t.Errorf("Expected score >= 80.0 for strong buy, got %f", score)
	}
}

func TestScoreCalculatorDifferentRatings(t *testing.T) {
	// Test different ratings produce different scores
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	baseTime := time.Now()

	// Strong buy stock
	strongBuyStock := []models.Stock{
		{
			Ticker:   "AAPL",
			RatingTo: "strong buy",
			TargetTo: "200.00",
			Time:     baseTime,
		},
	}

	// Sell stock
	sellStock := []models.Stock{
		{
			Ticker:   "AAPL",
			RatingTo: "sell",
			TargetTo: "100.00",
			Time:     baseTime,
		},
	}

	strongBuyScore := calculator.CalculateScore(strongBuyStock)
	sellScore := calculator.CalculateScore(sellStock)

	if strongBuyScore <= sellScore {
		t.Errorf("Expected strong buy score > sell score, got %f vs %f", strongBuyScore, sellScore)
	}
}

func TestScoreCalculatorWeightedCalculation(t *testing.T) {
	// Test that weighted calculation works properly
	weights := scoring.FactorWeights{
		Rating:      0.5,
		Target:      0.2,
		TargetTrend: 0.1,
		Volatility:  0.1,
		Recency:     0.1,
	}

	calculator := NewScoreCalculator(weights)

	stocks := []models.Stock{
		{
			Ticker:   "TEST",
			RatingTo: "buy",
			TargetTo: "150.00",
			Time:     time.Now(),
		},
	}

	score := calculator.CalculateScore(stocks)

	// Score should be positive and reasonable
	if score <= 0.0 {
		t.Errorf("Expected score > 0.0, got %f", score)
	}
	if score > 100.0 {
		t.Errorf("Expected score <= 100.0, got %f", score)
	}
}

func TestScoreCalculatorEdgeCases(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	// Test with invalid target prices
	stocks := []models.Stock{
		{
			Ticker:   "TEST",
			RatingTo: "buy",
			TargetTo: "invalid",
			Time:     time.Now(),
		},
	}

	score := calculator.CalculateScore(stocks)
	if score < 0.0 || score > 100.0 {
		t.Errorf("Expected score between 0.0 and 100.0, got %f", score)
	}

	// Test with empty ratings
	stocks[0].RatingTo = ""
	score = calculator.CalculateScore(stocks)
	if score < 0.0 || score > 100.0 {
		t.Errorf("Expected score between 0.0 and 100.0, got %f", score)
	}
}

func TestScoreCalculatorTimeOrdering(t *testing.T) {
	// Test that time ordering affects calculation
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	baseTime := time.Now()

	// Recent upgrade should score better than old downgrade
	stocks := []models.Stock{
		{
			Ticker:   "TEST",
			RatingTo: "buy",
			Time:     baseTime, // Recent
		},
		{
			Ticker:   "TEST",
			RatingTo: "sell",
			Time:     baseTime.Add(-30 * 24 * time.Hour), // Old
		},
	}

	score := calculator.CalculateScore(stocks)
	if score <= 50.0 {
		t.Errorf("Expected score > 50.0 due to recent upgrade, got %f", score)
	}
}

// Test para ratings desconocidos y logging
func TestScoreCalculatorUnknownRatings(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	// Crear el directorio logs si no existe
	if err := os.MkdirAll("logs", 0755); err != nil {
		t.Logf("Warning: Could not create logs directory: %v", err)
	}

	stocks := []models.Stock{
		{
			Ticker:     "TEST",
			RatingFrom: "unknown_rating_from",
			RatingTo:   "unknown_rating_to",
			TargetFrom: "100.00",
			TargetTo:   "150.00",
			Time:       time.Now(),
		},
	}

	score := calculator.CalculateScore(stocks)

	// Debe devolver un score válido incluso con ratings desconocidos
	if score < 0.0 || score > 100.0 {
		t.Errorf("Expected score between 0.0 and 100.0 for unknown ratings, got %f", score)
	}

	// Limpiar archivo de log si existe
	_ = os.Remove("logs/unknown_ratings.log")
}

// Test para extractNumber function
func TestExtractNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"$150.00", 150.0, false},
		{"1,250.50", 1250.5, false},
		{"  $2,000  ", 2000.0, false},
		{"invalid", 0.0, true},
		{"", 0.0, true},
		{"100", 100.0, false},
	}

	for _, test := range tests {
		result, err := extractNumber(test.input)

		if test.hasError {
			if err == nil {
				t.Errorf("Expected error for input '%s', got nil", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("Expected %f for input '%s', got %f", test.expected, test.input, result)
			}
		}
	}
}

// Test para calculateTargetScore con casos específicos
func TestCalculateTargetScore(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	tests := []struct {
		name       string
		targetFrom string
		targetTo   string
		expected   float64
	}{
		{"Positive change", "100.00", "120.00", 1.0},        // 20% increase
		{"Negative change", "100.00", "80.00", 0.0},         // 20% decrease
		{"Small positive change", "100.00", "105.00", 0.25}, // 5% increase
		{"Invalid from", "invalid", "120.00", 0.0},
		{"Invalid to", "100.00", "invalid", 0.0},
		{"Zero values", "0.00", "120.00", 0.0},
		{"Negative values", "-100.00", "120.00", 0.0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stock := models.Stock{
				TargetFrom: test.targetFrom,
				TargetTo:   test.targetTo,
			}

			score := calculator.calculateTargetScore(stock)
			if score != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, score)
			}
		})
	}
}

// Test para calculateTargetTrendScore con casos específicos
func TestCalculateTargetTrendScore(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	baseTime := time.Now()

	tests := []struct {
		name     string
		stocks   []models.Stock
		expected float64
	}{
		{
			name:     "Empty history",
			stocks:   []models.Stock{},
			expected: 0.5,
		},
		{
			name: "Single stock",
			stocks: []models.Stock{
				{TargetFrom: "100.00", TargetTo: "120.00", Time: baseTime},
			},
			expected: 0.5,
		},
		{
			name: "Consistent positive trend",
			stocks: []models.Stock{
				{TargetFrom: "100.00", TargetTo: "110.00", Time: baseTime},
				{TargetFrom: "90.00", TargetTo: "100.00", Time: baseTime.Add(-time.Hour)},
				{TargetFrom: "80.00", TargetTo: "90.00", Time: baseTime.Add(-2 * time.Hour)},
			},
			// Esta función calcula un score basado en tendencias y consistencia
			// No necesariamente será 1.0, verificamos que esté en un rango razonable
		},
		{
			name: "Old data beyond 6 months",
			stocks: []models.Stock{
				{TargetFrom: "100.00", TargetTo: "110.00", Time: baseTime.AddDate(0, -7, 0)},
			},
			expected: 0.5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score := calculator.calculateTargetTrendScore(test.stocks)

			if test.expected > 0 && score != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, score)
			} else if test.expected == 0 {
				// Para casos donde no esperamos un valor exacto, verificamos rango válido
				if score < 0.0 || score > 1.0 {
					t.Errorf("Score should be between 0.0 and 1.0, got %f", score)
				}
				// Para trend positivo, esperamos score > 0.5
				if test.name == "Consistent positive trend" && score <= 0.5 {
					t.Errorf("Expected score > 0.5 for positive trend, got %f", score)
				}
			}
		})
	}
}

// Test para calculateVolatilityScore
func TestCalculateVolatilityScore(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	baseTime := time.Now()

	tests := []struct {
		name     string
		stocks   []models.Stock
		expected float64
	}{
		{
			name:     "Less than 3 stocks",
			stocks:   []models.Stock{{RatingTo: "buy", Time: baseTime}},
			expected: 0.5,
		},
		{
			name: "Consistent ratings",
			stocks: []models.Stock{
				{RatingTo: "buy", Time: baseTime},
				{RatingTo: "buy", Time: baseTime.Add(-time.Hour)},
				{RatingTo: "buy", Time: baseTime.Add(-2 * time.Hour)},
			},
		},
		{
			name: "Old data beyond 3 months",
			stocks: []models.Stock{
				{RatingTo: "buy", Time: baseTime.AddDate(0, -4, 0)},
				{RatingTo: "buy", Time: baseTime.AddDate(0, -5, 0)},
				{RatingTo: "buy", Time: baseTime.AddDate(0, -6, 0)},
			},
			expected: 0.5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score := calculator.calculateVolatilityScore(test.stocks)
			if test.expected > 0 && score != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, score)
			}
			// Verificar que el score esté en el rango válido
			if score < 0.0 || score > 1.0 {
				t.Errorf("Score should be between 0.0 and 1.0, got %f", score)
			}
		})
	}
}

// Test para calculateRecencyScore
func TestCalculateRecencyScore(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		minScore float64
		maxScore float64
	}{
		{"Very recent", now, 0.9, 1.0},                      // Today should be close to 1.0
		{"30 days ago", now.AddDate(0, 0, -30), 0.4, 0.6},   // 30 days should be around 0.5
		{"90 days ago", now.AddDate(0, 0, -90), 0.15, 0.35}, // 90 days should be lower
		{"Very old", now.AddDate(-1, 0, 0), 0.0, 0.1},       // 1 year should be very low
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			score := calculator.calculateRecencyScore(test.time)

			// Verificar que el score esté en el rango válido
			if score < 0.0 || score > 1.0 {
				t.Errorf("Score should be between 0.0 and 1.0, got %f", score)
			}

			// Verificar que esté en el rango esperado
			if score < test.minScore || score > test.maxScore {
				t.Errorf("Expected score between %f and %f, got %f", test.minScore, test.maxScore, score)
			}
		})
	}
}

// Test para calculateRatingTrend con múltiples escenarios
func TestCalculateRatingTrend(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	baseTime := time.Now()

	tests := []struct {
		name     string
		stocks   []models.Stock
		testType string
	}{
		{
			name: "Insufficient data for trend",
			stocks: []models.Stock{
				{RatingTo: "buy", Time: baseTime},
			},
			testType: "immediate",
		},
		{
			name: "Consistent upgrade trend",
			stocks: []models.Stock{
				{RatingTo: "strong buy", Time: baseTime},
				{RatingTo: "buy", Time: baseTime.Add(-time.Hour)},
				{RatingTo: "hold", Time: baseTime.Add(-2 * time.Hour)},
				{RatingTo: "sell", Time: baseTime.Add(-3 * time.Hour)},
			},
			testType: "trend",
		},
		{
			name: "Old data beyond 6 months",
			stocks: []models.Stock{
				{RatingTo: "buy", Time: baseTime.AddDate(0, -7, 0)},
				{RatingTo: "hold", Time: baseTime.AddDate(0, -8, 0)},
				{RatingTo: "sell", Time: baseTime.AddDate(0, -9, 0)},
			},
			testType: "immediate",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bonus := calculator.calculateRatingTrend(test.stocks)

			// Verificar que el bonus esté en el rango válido
			if bonus < -0.15 || bonus > 0.15 {
				t.Errorf("Bonus should be between -0.15 and 0.15, got %f", bonus)
			}
		})
	}
}

// Test para calculateImmediateChange
func TestCalculateImmediateChange(t *testing.T) {
	calculator := NewScoreCalculator(scoring.DefaultWeights)

	tests := []struct {
		name     string
		stocks   []models.Stock
		expected float64
	}{
		{
			name:     "No stocks",
			stocks:   []models.Stock{},
			expected: 0.0,
		},
		{
			name: "Upgrade from sell to buy",
			stocks: []models.Stock{
				{RatingFrom: "sell", RatingTo: "buy", Time: time.Now()},
			},
		},
		{
			name: "Downgrade from buy to sell",
			stocks: []models.Stock{
				{RatingFrom: "buy", RatingTo: "sell", Time: time.Now()},
			},
		},
		{
			name: "Unknown ratings",
			stocks: []models.Stock{
				{RatingFrom: "unknown_from", RatingTo: "unknown_to", Time: time.Now()},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			change := calculator.calculateImmediateChange(test.stocks)

			if test.expected != 0 && change != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, change)
			}

			// Verificar que el cambio esté en el rango válido
			if change < -0.15 || change > 0.15 {
				t.Errorf("Change should be between -0.15 and 0.15, got %f", change)
			}
		})
	}
}
