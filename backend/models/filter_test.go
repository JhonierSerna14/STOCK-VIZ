package models

import (
	"testing"
	"time"
)

func TestRecommendationFilterCreation(t *testing.T) {
	// Test creating a new RecommendationFilter instance
	filter := RecommendationFilter{
		Limit:    10,
		DateFrom: "2024-01-01",
		DateTo:   "2024-12-31",
		Rating:   "buy",
		Ticker:   "AAPL",
	}

	if filter.Limit != 10 {
		t.Errorf("Expected Limit 10, got %d", filter.Limit)
	}
	if filter.DateFrom != "2024-01-01" {
		t.Errorf("Expected DateFrom 2024-01-01, got %s", filter.DateFrom)
	}
	if filter.DateTo != "2024-12-31" {
		t.Errorf("Expected DateTo 2024-12-31, got %s", filter.DateTo)
	}
	if filter.Rating != "buy" {
		t.Errorf("Expected Rating buy, got %s", filter.Rating)
	}
	if filter.Ticker != "AAPL" {
		t.Errorf("Expected Ticker AAPL, got %s", filter.Ticker)
	}
}

func TestRecommendationFilterEmptyValues(t *testing.T) {
	// Test RecommendationFilter with empty values
	filter := RecommendationFilter{}

	if filter.Limit != 0 {
		t.Errorf("Expected Limit 0, got %d", filter.Limit)
	}
	if filter.DateFrom != "" {
		t.Errorf("Expected empty DateFrom, got %s", filter.DateFrom)
	}
	if filter.DateTo != "" {
		t.Errorf("Expected empty DateTo, got %s", filter.DateTo)
	}
	if filter.Rating != "" {
		t.Errorf("Expected empty Rating, got %s", filter.Rating)
	}
	if filter.Ticker != "" {
		t.Errorf("Expected empty Ticker, got %s", filter.Ticker)
	}
}

func TestRecommendationFilterParseDatesValidDates(t *testing.T) {
	// Test parsing valid dates
	filter := RecommendationFilter{
		DateFrom: "2024-01-01",
		DateTo:   "2024-12-31",
	}

	from, to, err := filter.ParseDates()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if from == nil {
		t.Fatal("Expected from date, got nil")
	}
	if to == nil {
		t.Fatal("Expected to date, got nil")
	}

	if from.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", from.Year())
	}
	if from.Month() != time.January {
		t.Errorf("Expected month January, got %v", from.Month())
	}
	if from.Day() != 1 {
		t.Errorf("Expected day 1, got %d", from.Day())
	}

	if to.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", to.Year())
	}
	if to.Month() != time.December {
		t.Errorf("Expected month December, got %v", to.Month())
	}
	if to.Day() != 31 {
		t.Errorf("Expected day 31, got %d", to.Day())
	}
	// Check that 'to' date is set to end of day
	if to.Hour() != 23 {
		t.Errorf("Expected hour 23, got %d", to.Hour())
	}
	if to.Minute() != 59 {
		t.Errorf("Expected minute 59, got %d", to.Minute())
	}
	if to.Second() != 59 {
		t.Errorf("Expected second 59, got %d", to.Second())
	}
}

func TestRecommendationFilterParseDatesEmptyDates(t *testing.T) {
	// Test parsing empty dates
	filter := RecommendationFilter{}

	from, to, err := filter.ParseDates()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if from != nil {
		t.Error("Expected nil from date")
	}
	if to != nil {
		t.Error("Expected nil to date")
	}
}

func TestRecommendationFilterParseDatesOnlyFromDate(t *testing.T) {
	// Test parsing only from date
	filter := RecommendationFilter{
		DateFrom: "2024-06-15",
	}

	from, to, err := filter.ParseDates()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if from == nil {
		t.Fatal("Expected from date, got nil")
	}
	if to != nil {
		t.Error("Expected nil to date")
	}

	if from.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", from.Year())
	}
	if from.Month() != time.June {
		t.Errorf("Expected month June, got %v", from.Month())
	}
	if from.Day() != 15 {
		t.Errorf("Expected day 15, got %d", from.Day())
	}
}

func TestRecommendationFilterParseDatesOnlyToDate(t *testing.T) {
	// Test parsing only to date
	filter := RecommendationFilter{
		DateTo: "2024-06-15",
	}

	from, to, err := filter.ParseDates()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if from != nil {
		t.Error("Expected nil from date")
	}
	if to == nil {
		t.Fatal("Expected to date, got nil")
	}

	if to.Year() != 2024 {
		t.Errorf("Expected year 2024, got %d", to.Year())
	}
	if to.Month() != time.June {
		t.Errorf("Expected month June, got %v", to.Month())
	}
	if to.Day() != 15 {
		t.Errorf("Expected day 15, got %d", to.Day())
	}
}

func TestRecommendationFilterParseDatesInvalidDateFormat(t *testing.T) {
	tests := []struct {
		name     string
		dateFrom string
		dateTo   string
	}{
		{"Invalid from date", "invalid-date", ""},
		{"Invalid to date", "", "invalid-date"},
		{"Both invalid", "invalid-from", "invalid-to"},
		{"Wrong format from", "01/01/2024", ""},
		{"Wrong format to", "", "31/12/2024"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := RecommendationFilter{
				DateFrom: tt.dateFrom,
				DateTo:   tt.dateTo,
			}

			_, _, err := filter.ParseDates()
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestRecommendationFilterParseDatesEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		dateFrom string
		dateTo   string
	}{
		{"Leap year", "2024-02-29", "2024-02-29"},
		{"End of month", "2024-01-31", "2024-01-31"},
		{"Start of year", "2024-01-01", "2024-01-01"},
		{"End of year", "2024-12-31", "2024-12-31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := RecommendationFilter{
				DateFrom: tt.dateFrom,
				DateTo:   tt.dateTo,
			}

			from, to, err := filter.ParseDates()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if from == nil {
				t.Fatal("Expected from date, got nil")
			}
			if to == nil {
				t.Fatal("Expected to date, got nil")
			}
		})
	}
}
