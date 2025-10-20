// Package testhelpers provides common testing utilities and helpers for the STOCK-VIZ project
package testhelpers

import (
	"testing"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	// Use a test PostgreSQL database or skip if not available
	// For now, we'll skip SQLite tests and use mocks instead
	t.Skip("SQLite tests require CGO. Use integration tests with PostgreSQL for database testing.")
	return nil
}

// CreateTestStock creates a test stock with default values
func CreateTestStock(ticker, company string, timeOffset time.Duration) models.Stock {
	return models.Stock{
		Ticker:     ticker,
		Company:    company,
		Brokerage:  "Test Brokerage",
		Action:     "upgrade",
		RatingFrom: "hold",
		RatingTo:   "buy",
		TargetFrom: "100.00",
		TargetTo:   "150.00",
		Time:       time.Now().Add(timeOffset),
	}
}

// CreateTestStockWithRating creates a test stock with specific rating
func CreateTestStockWithRating(ticker, company, rating string, timeOffset time.Duration) models.Stock {
	return models.Stock{
		Ticker:     ticker,
		Company:    company,
		Brokerage:  "Test Brokerage",
		Action:     "upgrade",
		RatingFrom: "hold",
		RatingTo:   rating,
		TargetFrom: "100.00",
		TargetTo:   "150.00",
		Time:       time.Now().Add(timeOffset),
	}
}

// CreateTestStockFull creates a test stock with all parameters customizable
func CreateTestStockFull(ticker, company, brokerage, action, ratingFrom, ratingTo, targetFrom, targetTo string, timeOffset time.Duration) models.Stock {
	return models.Stock{
		Ticker:     ticker,
		Company:    company,
		Brokerage:  brokerage,
		Action:     action,
		RatingFrom: ratingFrom,
		RatingTo:   ratingTo,
		TargetFrom: targetFrom,
		TargetTo:   targetTo,
		Time:       time.Now().Add(timeOffset),
	}
}

// SeedTestData seeds the test database with sample data
func SeedTestData(t *testing.T, db *gorm.DB) {
	t.Helper()

	stocks := []models.Stock{
		CreateTestStock("AAPL", "Apple Inc", 0),
		CreateTestStock("GOOGL", "Alphabet Inc", -time.Hour),
		CreateTestStock("MSFT", "Microsoft Corp", -2*time.Hour),
		CreateTestStockWithRating("TSLA", "Tesla Inc", "strong buy", -time.Hour),
		CreateTestStockWithRating("AMZN", "Amazon Inc", "sell", -3*time.Hour),
	}

	for _, stock := range stocks {
		if err := db.Create(&stock).Error; err != nil {
			t.Fatalf("Failed to seed test data: %v", err)
		}
	}
}

// AssertNoError is a helper to assert no error occurred
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError is a helper to assert an error occurred
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// AssertEqual is a generic helper for equality assertions
func AssertEqual[T comparable](t *testing.T, expected, actual T) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

// AssertNotEmpty asserts that a string is not empty
func AssertNotEmpty(t *testing.T, value string) {
	t.Helper()
	if value == "" {
		t.Error("Expected non-empty string")
	}
}

// AssertSliceLength asserts the length of a slice
func AssertSliceLength[T any](t *testing.T, slice []T, expectedLength int) {
	t.Helper()
	if len(slice) != expectedLength {
		t.Errorf("Expected slice length %d, got %d", expectedLength, len(slice))
	}
}

// AssertInRange asserts that a number is within a specific range
func AssertInRange(t *testing.T, value, min, max float64) {
	t.Helper()
	if value < min || value > max {
		t.Errorf("Expected value between %f and %f, got %f", min, max, value)
	}
}
