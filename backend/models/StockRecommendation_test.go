package models

import (
	"testing"
	"time"
)

func TestStockRecommendationCreation(t *testing.T) {
	// Test creating a new StockRecommendation instance
	testTime := time.Now()
	recommendation := StockRecommendation{
		Ticker:            "AAPL",
		Company:           "Apple Inc",
		Score:             85.5,
		LatestRating:      "buy",
		LatestTarget:      "175.00",
		LastUpdated:       testTime,
		AnalysisRationale: "Strong fundamentals and positive outlook",
	}

	if recommendation.Ticker != "AAPL" {
		t.Errorf("Expected Ticker AAPL, got %s", recommendation.Ticker)
	}
	if recommendation.Company != "Apple Inc" {
		t.Errorf("Expected Company Apple Inc, got %s", recommendation.Company)
	}
	if recommendation.Score != 85.5 {
		t.Errorf("Expected Score 85.5, got %f", recommendation.Score)
	}
	if recommendation.LatestRating != "buy" {
		t.Errorf("Expected LatestRating buy, got %s", recommendation.LatestRating)
	}
	if recommendation.LatestTarget != "175.00" {
		t.Errorf("Expected LatestTarget 175.00, got %s", recommendation.LatestTarget)
	}
	if recommendation.LastUpdated != testTime {
		t.Errorf("Expected LastUpdated %v, got %v", testTime, recommendation.LastUpdated)
	}
	if recommendation.AnalysisRationale != "Strong fundamentals and positive outlook" {
		t.Errorf("Expected AnalysisRationale 'Strong fundamentals and positive outlook', got %s", recommendation.AnalysisRationale)
	}
}

func TestStockRecommendationEmptyValues(t *testing.T) {
	// Test StockRecommendation with empty values
	recommendation := StockRecommendation{}

	if recommendation.Ticker != "" {
		t.Errorf("Expected empty Ticker, got %s", recommendation.Ticker)
	}
	if recommendation.Company != "" {
		t.Errorf("Expected empty Company, got %s", recommendation.Company)
	}
	if recommendation.Score != 0.0 {
		t.Errorf("Expected Score 0.0, got %f", recommendation.Score)
	}
	if recommendation.LatestRating != "" {
		t.Errorf("Expected empty LatestRating, got %s", recommendation.LatestRating)
	}
	if recommendation.LatestTarget != "" {
		t.Errorf("Expected empty LatestTarget, got %s", recommendation.LatestTarget)
	}
	if !recommendation.LastUpdated.IsZero() {
		t.Error("Expected zero LastUpdated")
	}
	if recommendation.AnalysisRationale != "" {
		t.Errorf("Expected empty AnalysisRationale, got %s", recommendation.AnalysisRationale)
	}
}

func TestStockRecommendationScoreValidation(t *testing.T) {
	tests := []struct {
		name  string
		score float64
	}{
		{"Zero score", 0.0},
		{"Low score", 25.5},
		{"Medium score", 50.0},
		{"High score", 85.5},
		{"Maximum score", 100.0},
		{"Above maximum", 150.0}, // Edge case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendation := StockRecommendation{
				Ticker: "TEST",
				Score:  tt.score,
			}
			if recommendation.Score != tt.score {
				t.Errorf("Expected Score %f, got %f", tt.score, recommendation.Score)
			}
		})
	}
}
