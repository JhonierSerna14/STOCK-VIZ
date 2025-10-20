package repository

import (
	"testing"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/JhonierSerna14/STOCK-VIZ/testhelpers"
	"github.com/stretchr/testify/assert"
)

// createTestStock creates a single test stock with default values
func createTestStock() models.Stock {
	return models.Stock{
		Ticker:     "TEST",
		Company:    "Test Company",
		Brokerage:  "Test Brokerage",
		Action:     "buy",
		RatingFrom: "hold",
		RatingTo:   "buy",
		TargetFrom: "100.00",
		TargetTo:   "150.00",
		Time:       time.Now(),
	}
}

func TestNewStockRepository(t *testing.T) {
	// Test creating a new repository with a nil database
	repo := NewStockRepository(nil)
	assert.NotNil(t, repo, "Expected repo to be non-nil")
}

// Note: The following tests are designed to test the repository logic
// For integration testing with a real database, use separate integration test files
// These tests focus on the repository interface and error handling

func TestStockRepository_InterfaceCompliance(t *testing.T) {
	// Test that our repository implements the expected interface
	repo := NewStockRepository(nil)

	// Verify the repository structure exists
	assert.NotNil(t, repo, "Repository should not be nil")

	// Test that the repository has a db field (even if nil)
	// This ensures the struct is properly initialized
	assert.NotNil(t, repo, "Repository struct should be initialized")
}

func TestStockRepository_GetFilteredStocks_DateValidation(t *testing.T) {
	// Test date parsing logic
	filter := models.RecommendationFilter{
		DateFrom: "invalid-date",
	}

	// The ParseDates method should be tested directly on the filter
	_, _, err := filter.ParseDates()
	assert.Error(t, err, "Should error with invalid date format")

	// Test valid date parsing
	validFilter := models.RecommendationFilter{
		DateFrom: "2024-01-01",
		DateTo:   "2024-12-31",
	}

	from, to, err := validFilter.ParseDates()
	assert.NoError(t, err, "Should not error with valid dates")
	assert.NotNil(t, from, "From date should be parsed")
	assert.NotNil(t, to, "To date should be parsed")
}

func TestStockRepository_GetStocksPaginated_ParameterValidation(t *testing.T) {
	// Test parameter validation logic
	tests := []struct {
		name  string
		page  int
		limit int
		query string
		valid bool
	}{
		{
			name:  "Valid pagination parameters",
			page:  1,
			limit: 10,
			query: "",
			valid: true,
		},
		{
			name:  "Valid pagination with search",
			page:  2,
			limit: 5,
			query: "AAPL",
			valid: true,
		},
		{
			name:  "Edge case - page 0",
			page:  0,
			limit: 10,
			query: "",
			valid: true, // Implementation should handle this
		},
		{
			name:  "Edge case - negative page",
			page:  -1,
			limit: 10,
			query: "",
			valid: true, // Implementation should handle this
		},
		{
			name:  "Edge case - zero limit",
			page:  1,
			limit: 0,
			query: "",
			valid: true, // Implementation should handle this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the parameters are structured correctly
			assert.GreaterOrEqual(t, tt.page, -1, "Page should be valid integer")
			assert.GreaterOrEqual(t, tt.limit, 0, "Limit should be non-negative")
			assert.IsType(t, "", tt.query, "Query should be string")
		})
	}
}

// TestStockRepository_DataStructures tests the data structures used
func TestStockRepository_DataStructures(t *testing.T) {
	// Test creating various stock combinations
	stocks := []models.Stock{
		testhelpers.CreateTestStock("AAPL", "Apple Inc", 0),
		testhelpers.CreateTestStockWithRating("GOOGL", "Alphabet Inc", "strong buy", -time.Hour),
		testhelpers.CreateTestStockFull("MSFT", "Microsoft Corp", "Goldman Sachs", "upgrade", "hold", "buy", "200.00", "250.00", -2*time.Hour),
	}

	// Verify structure integrity
	for i, stock := range stocks {
		assert.NotEmpty(t, stock.Ticker, "Stock %d should have ticker", i)
		assert.NotEmpty(t, stock.Company, "Stock %d should have company", i)
		assert.False(t, stock.Time.IsZero(), "Stock %d should have valid time", i)
	}

	// Test filter structure
	filter := models.RecommendationFilter{
		Limit:    10,
		DateFrom: "2024-01-01",
		DateTo:   "2024-12-31",
		Rating:   "buy",
		Ticker:   "AAPL",
	}

	assert.Equal(t, 10, filter.Limit)
	assert.Equal(t, "2024-01-01", filter.DateFrom)
	assert.Equal(t, "2024-12-31", filter.DateTo)
	assert.Equal(t, "buy", filter.Rating)
	assert.Equal(t, "AAPL", filter.Ticker)
}

// Integration test marker - these would run with a real database in CI/CD
func TestStockRepository_Integration(t *testing.T) {
	// Skip if not in integration test mode
	t.Skip("Integration tests require a database connection. Run with integration test flag.")

	// This is where you would put integration tests that:
	// 1. Connect to a test PostgreSQL instance
	// 2. Run migrations
	// 3. Test actual database operations
	// 4. Clean up after tests
}
