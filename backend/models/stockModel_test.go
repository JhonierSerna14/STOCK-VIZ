package models

import (
	"testing"
	"time"
)

func TestStockCreation(t *testing.T) {
	// Test creating a new Stock instance
	testTime := time.Now()
	stock := Stock{
		Ticker:     "AAPL",
		Company:    "Apple Inc",
		Brokerage:  "Goldman Sachs",
		Action:     "upgrade",
		RatingFrom: "neutral",
		RatingTo:   "buy",
		TargetFrom: "150.00",
		TargetTo:   "175.00",
		Time:       testTime,
	}

	if stock.Ticker != "AAPL" {
		t.Errorf("Expected Ticker AAPL, got %s", stock.Ticker)
	}
	if stock.Company != "Apple Inc" {
		t.Errorf("Expected Company Apple Inc, got %s", stock.Company)
	}
	if stock.Brokerage != "Goldman Sachs" {
		t.Errorf("Expected Brokerage Goldman Sachs, got %s", stock.Brokerage)
	}
	if stock.Action != "upgrade" {
		t.Errorf("Expected Action upgrade, got %s", stock.Action)
	}
	if stock.RatingFrom != "neutral" {
		t.Errorf("Expected RatingFrom neutral, got %s", stock.RatingFrom)
	}
	if stock.RatingTo != "buy" {
		t.Errorf("Expected RatingTo buy, got %s", stock.RatingTo)
	}
	if stock.TargetFrom != "150.00" {
		t.Errorf("Expected TargetFrom 150.00, got %s", stock.TargetFrom)
	}
	if stock.TargetTo != "175.00" {
		t.Errorf("Expected TargetTo 175.00, got %s", stock.TargetTo)
	}
	if stock.Time != testTime {
		t.Errorf("Expected Time %v, got %v", testTime, stock.Time)
	}
}

func TestStockEmptyValues(t *testing.T) {
	// Test Stock with empty values
	stock := Stock{}

	if stock.Ticker != "" {
		t.Errorf("Expected empty Ticker, got %s", stock.Ticker)
	}
	if stock.Company != "" {
		t.Errorf("Expected empty Company, got %s", stock.Company)
	}
	if stock.Brokerage != "" {
		t.Errorf("Expected empty Brokerage, got %s", stock.Brokerage)
	}
	if stock.Action != "" {
		t.Errorf("Expected empty Action, got %s", stock.Action)
	}
	if stock.RatingFrom != "" {
		t.Errorf("Expected empty RatingFrom, got %s", stock.RatingFrom)
	}
	if stock.RatingTo != "" {
		t.Errorf("Expected empty RatingTo, got %s", stock.RatingTo)
	}
	if stock.TargetFrom != "" {
		t.Errorf("Expected empty TargetFrom, got %s", stock.TargetFrom)
	}
	if stock.TargetTo != "" {
		t.Errorf("Expected empty TargetTo, got %s", stock.TargetTo)
	}
	if !stock.Time.IsZero() {
		t.Error("Expected zero Time")
	}
}

func TestStockPrimaryKeyFields(t *testing.T) {
	// Test that primary key fields are properly set
	testTime := time.Now()
	stock := Stock{
		Ticker: "MSFT",
		Time:   testTime,
	}

	if stock.Ticker != "MSFT" {
		t.Errorf("Expected Ticker MSFT, got %s", stock.Ticker)
	}
	if stock.Time != testTime {
		t.Errorf("Expected Time %v, got %v", testTime, stock.Time)
	}
}
