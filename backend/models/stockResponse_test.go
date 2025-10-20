package models

import (
	"testing"
	"time"
)

func TestStockResponseCreation(t *testing.T) {
	// Test creating a new StockResponse instance
	testTime := time.Now()
	stocks := []Stock{
		{
			Ticker:  "AAPL",
			Company: "Apple Inc",
			Time:    testTime,
		},
		{
			Ticker:  "GOOGL",
			Company: "Alphabet Inc",
			Time:    testTime.Add(-time.Hour),
		},
	}

	response := StockResponse{
		Items:    stocks,
		NextPage: "page_token_123",
	}

	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}
	if response.Items[0].Ticker != "AAPL" {
		t.Errorf("Expected first ticker AAPL, got %s", response.Items[0].Ticker)
	}
	if response.Items[1].Ticker != "GOOGL" {
		t.Errorf("Expected second ticker GOOGL, got %s", response.Items[1].Ticker)
	}
	if response.NextPage != "page_token_123" {
		t.Errorf("Expected NextPage page_token_123, got %s", response.NextPage)
	}
}

func TestStockResponseEmptyValues(t *testing.T) {
	// Test StockResponse with empty values
	response := StockResponse{}

	if response.Items != nil {
		t.Error("Expected nil Items")
	}
	if response.NextPage != "" {
		t.Errorf("Expected empty NextPage, got %s", response.NextPage)
	}
}

func TestStockResponseEmptyItems(t *testing.T) {
	// Test StockResponse with empty items slice
	response := StockResponse{
		Items:    []Stock{},
		NextPage: "",
	}

	if response.Items == nil {
		t.Error("Expected non-nil Items")
	}
	if len(response.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(response.Items))
	}
	if response.NextPage != "" {
		t.Errorf("Expected empty NextPage, got %s", response.NextPage)
	}
}

func TestStockResponseNoNextPage(t *testing.T) {
	// Test StockResponse without next page
	stocks := []Stock{
		{Ticker: "AAPL", Company: "Apple Inc"},
	}

	response := StockResponse{
		Items:    stocks,
		NextPage: "",
	}

	if len(response.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(response.Items))
	}
	if response.NextPage != "" {
		t.Errorf("Expected empty NextPage, got %s", response.NextPage)
	}
}
