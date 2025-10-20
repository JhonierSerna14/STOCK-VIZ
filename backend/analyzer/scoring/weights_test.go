package scoring

import (
	"math"
	"testing"
)

func TestFactorWeightsDefaultWeights(t *testing.T) {
	// Test that default weights are properly defined
	weights := DefaultWeights

	if weights.Rating != 0.30 {
		t.Errorf("Expected Rating 0.30, got %f", weights.Rating)
	}
	if weights.Target != 0.25 {
		t.Errorf("Expected Target 0.25, got %f", weights.Target)
	}
	if weights.TargetTrend != 0.20 {
		t.Errorf("Expected TargetTrend 0.20, got %f", weights.TargetTrend)
	}
	if weights.Volatility != 0.15 {
		t.Errorf("Expected Volatility 0.15, got %f", weights.Volatility)
	}
	if weights.Recency != 0.10 {
		t.Errorf("Expected Recency 0.10, got %f", weights.Recency)
	}
}

func TestFactorWeightsSumToOne(t *testing.T) {
	// Test that weights sum to approximately 1.0
	weights := DefaultWeights
	sum := weights.Rating + weights.Target + weights.TargetTrend + weights.Volatility + weights.Recency

	delta := 0.001
	if math.Abs(sum-1.0) > delta {
		t.Errorf("Expected sum to be approximately 1.0 (within %f), got %f", delta, sum)
	}
}

func TestFactorWeightsValidRange(t *testing.T) {
	// Test that all weights are between 0 and 1
	weights := DefaultWeights

	weightsToTest := map[string]float64{
		"Rating":      weights.Rating,
		"Target":      weights.Target,
		"TargetTrend": weights.TargetTrend,
		"Volatility":  weights.Volatility,
		"Recency":     weights.Recency,
	}

	for name, weight := range weightsToTest {
		if weight < 0.0 {
			t.Errorf("Weight %s should be >= 0.0, got %f", name, weight)
		}
		if weight > 1.0 {
			t.Errorf("Weight %s should be <= 1.0, got %f", name, weight)
		}
	}
}

func TestFactorWeightsCustomWeights(t *testing.T) {
	// Test creating custom weights
	customWeights := FactorWeights{
		Rating:      0.40,
		Target:      0.30,
		TargetTrend: 0.15,
		Volatility:  0.10,
		Recency:     0.05,
	}

	if customWeights.Rating != 0.40 {
		t.Errorf("Expected Rating 0.40, got %f", customWeights.Rating)
	}
	if customWeights.Target != 0.30 {
		t.Errorf("Expected Target 0.30, got %f", customWeights.Target)
	}
	if customWeights.TargetTrend != 0.15 {
		t.Errorf("Expected TargetTrend 0.15, got %f", customWeights.TargetTrend)
	}
	if customWeights.Volatility != 0.10 {
		t.Errorf("Expected Volatility 0.10, got %f", customWeights.Volatility)
	}
	if customWeights.Recency != 0.05 {
		t.Errorf("Expected Recency 0.05, got %f", customWeights.Recency)
	}

	// Verify custom weights sum to 1.0
	sum := customWeights.Rating + customWeights.Target + customWeights.TargetTrend + customWeights.Volatility + customWeights.Recency
	delta := 0.001
	if math.Abs(sum-1.0) > delta {
		t.Errorf("Expected sum to be approximately 1.0 (within %f), got %f", delta, sum)
	}
}

func TestFactorWeightsZeroValues(t *testing.T) {
	// Test with zero weights (edge case)
	zeroWeights := FactorWeights{
		Rating:      0.0,
		Target:      0.0,
		TargetTrend: 0.0,
		Volatility:  0.0,
		Recency:     1.0,
	}

	if zeroWeights.Rating != 0.0 {
		t.Errorf("Expected Rating 0.0, got %f", zeroWeights.Rating)
	}
	if zeroWeights.Target != 0.0 {
		t.Errorf("Expected Target 0.0, got %f", zeroWeights.Target)
	}
	if zeroWeights.TargetTrend != 0.0 {
		t.Errorf("Expected TargetTrend 0.0, got %f", zeroWeights.TargetTrend)
	}
	if zeroWeights.Volatility != 0.0 {
		t.Errorf("Expected Volatility 0.0, got %f", zeroWeights.Volatility)
	}
	if zeroWeights.Recency != 1.0 {
		t.Errorf("Expected Recency 1.0, got %f", zeroWeights.Recency)
	}
}
