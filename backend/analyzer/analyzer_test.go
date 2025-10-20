package analyzer

import (
	"testing"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/mocks"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/JhonierSerna14/STOCK-VIZ/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestData helpers
func createTestStock(ticker, company, ratingTo, targetTo string, timeOffset time.Duration) models.Stock {
	return models.Stock{
		Ticker:   ticker,
		Company:  company,
		RatingTo: ratingTo,
		TargetTo: targetTo,
		Time:     time.Now().Add(timeOffset),
	}
}

func TestNewStockAnalyzer(t *testing.T) {
	// Test creating analyzer with mock repository
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	assert.NotNil(t, analyzer, "Analyzer should not be nil")
	assert.NotNil(t, analyzer.repository, "Repository should be set")
	assert.NotNil(t, analyzer.calculator, "Calculator should be initialized")
	assert.NotNil(t, analyzer.formatter, "Formatter should be initialized")
}

func TestStockAnalyzer_GetFilteredRecommendations(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	// Setup test data
	baseTime := time.Now()
	allStocks := []models.Stock{
		testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Goldman Sachs", "upgrade", "hold", "buy", "150.00", "175.00", 0),
		testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Morgan Stanley", "upgrade", "neutral", "hold", "140.00", "150.00", -24*time.Hour),
		testhelpers.CreateTestStockFull("GOOGL", "Alphabet Inc", "JP Morgan", "upgrade", "hold", "strong buy", "2500.00", "2800.00", 0),
		testhelpers.CreateTestStockFull("MSFT", "Microsoft Corp", "Goldman Sachs", "downgrade", "buy", "hold", "300.00", "280.00", 0),
	}

	// Set specific times for predictable ordering
	for i := range allStocks {
		allStocks[i].Time = baseTime.Add(time.Duration(-i) * time.Hour)
	}

	// Setup mock for empty filter - return all stocks
	mockRepo.On("GetFilteredStocks", models.RecommendationFilter{}).Return(allStocks, nil)

	// Setup mock for AAPL filter - return only AAPL stocks
	aaplStocks := []models.Stock{allStocks[0], allStocks[1]} // First two stocks are AAPL
	mockRepo.On("GetFilteredStocks", models.RecommendationFilter{Ticker: "AAPL"}).Return(aaplStocks, nil)

	// Setup mock for limit filter - return all stocks (analyzer will limit)
	mockRepo.On("GetFilteredStocks", models.RecommendationFilter{Limit: 2}).Return(allStocks, nil)

	// Setup mock for rating filter - return AAPL stocks only (they have "buy" rating)
	mockRepo.On("GetFilteredStocks", models.RecommendationFilter{Rating: "buy"}).Return(aaplStocks, nil)

	// Setup mock for nonexistent ticker - return empty
	mockRepo.On("GetFilteredStocks", models.RecommendationFilter{Ticker: "NONEXISTENT"}).Return([]models.Stock{}, nil)

	tests := []struct {
		name                string
		filter              models.RecommendationFilter
		expectedCount       int
		expectedFirstTicker string
		shouldError         bool
	}{
		{
			name:                "Get all recommendations",
			filter:              models.RecommendationFilter{},
			expectedCount:       3,  // AAPL, GOOGL, MSFT grouped by ticker
			expectedFirstTicker: "", // Order depends on scoring
			shouldError:         false,
		},
		{
			name: "Filter by ticker",
			filter: models.RecommendationFilter{
				Ticker: "AAPL",
			},
			expectedCount:       1,
			expectedFirstTicker: "AAPL",
			shouldError:         false,
		},
		{
			name: "Filter with limit",
			filter: models.RecommendationFilter{
				Limit: 2,
			},
			expectedCount: 2,
			shouldError:   false,
		},
		{
			name: "Filter by rating",
			filter: models.RecommendationFilter{
				Rating: "buy",
			},
			expectedCount: 1, // Only AAPL has latest rating "buy"
			shouldError:   false,
		},
		{
			name: "Empty result filter",
			filter: models.RecommendationFilter{
				Ticker: "NONEXISTENT",
			},
			expectedCount: 0,
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendations, err := analyzer.GetFilteredRecommendations(tt.filter)

			if tt.shouldError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, recommendations, tt.expectedCount)

			if tt.expectedFirstTicker != "" && len(recommendations) > 0 {
				assert.Equal(t, tt.expectedFirstTicker, recommendations[0].Ticker)
			}

			// Verify recommendations are sorted by score (descending)
			for i := 1; i < len(recommendations); i++ {
				assert.GreaterOrEqual(t, recommendations[i-1].Score, recommendations[i].Score,
					"Recommendations should be sorted by score descending")
			}

			// Verify all recommendations have required fields
			for _, rec := range recommendations {
				assert.NotEmpty(t, rec.Ticker, "Ticker should not be empty")
				assert.NotEmpty(t, rec.Company, "Company should not be empty")
				assert.GreaterOrEqual(t, rec.Score, 0.0, "Score should be non-negative")
				assert.LessOrEqual(t, rec.Score, 100.0, "Score should not exceed 100")
				assert.NotEmpty(t, rec.AnalysisRationale, "Analysis rationale should not be empty")
				assert.False(t, rec.LastUpdated.IsZero(), "LastUpdated should be set")
			}
		})
	}
}

func TestStockAnalyzer_GetFilteredRecommendations_RepositoryError(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	// Set up repository to return error
	mockRepo.On("GetFilteredStocks", mock.AnythingOfType("models.RecommendationFilter")).Return([]models.Stock{}, assert.AnError)

	filter := models.RecommendationFilter{}
	recommendations, err := analyzer.GetFilteredRecommendations(filter)

	assert.Error(t, err)
	assert.Nil(t, recommendations)
}

func TestStockAnalyzer_generateRecommendations(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	baseTime := time.Now()
	stockMap := map[string][]models.Stock{
		"AAPL": {
			testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Goldman", "upgrade", "hold", "buy", "150.00", "175.00", 0),
			testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Morgan", "upgrade", "neutral", "hold", "140.00", "150.00", -24*time.Hour),
		},
		"GOOGL": {
			testhelpers.CreateTestStockFull("GOOGL", "Alphabet Inc", "JP Morgan", "upgrade", "hold", "strong buy", "2500.00", "2800.00", 0),
		},
	}

	// Set specific times
	for ticker := range stockMap {
		for i := range stockMap[ticker] {
			stockMap[ticker][i].Time = baseTime.Add(time.Duration(-i) * time.Hour)
		}
	}

	recommendations := analyzer.generateRecommendations(stockMap)

	assert.Len(t, recommendations, 2, "Should generate one recommendation per ticker")

	// Verify recommendations structure
	for _, rec := range recommendations {
		assert.NotEmpty(t, rec.Ticker, "Ticker should not be empty")
		assert.NotEmpty(t, rec.Company, "Company should not be empty")
		assert.GreaterOrEqual(t, rec.Score, 0.0, "Score should be non-negative")
		assert.LessOrEqual(t, rec.Score, 100.0, "Score should not exceed 100")
		assert.NotEmpty(t, rec.AnalysisRationale, "Analysis rationale should not be empty")
		assert.False(t, rec.LastUpdated.IsZero(), "LastUpdated should be set")

		// Verify latest data is used
		if rec.Ticker == "AAPL" {
			assert.Equal(t, "buy", rec.LatestRating, "Should use latest rating")
			assert.Equal(t, "175.00", rec.LatestTarget, "Should use latest target")
		} else if rec.Ticker == "GOOGL" {
			assert.Equal(t, "strong buy", rec.LatestRating, "Should use latest rating")
			assert.Equal(t, "2800.00", rec.LatestTarget, "Should use latest target")
		}
	}
}

func TestStockAnalyzer_generateRecommendations_EmptyData(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	// Test with empty stock map
	stockMap := map[string][]models.Stock{}
	recommendations := analyzer.generateRecommendations(stockMap)

	assert.Empty(t, recommendations, "Should return empty slice for empty data")
}

func TestStockAnalyzer_generateRationale(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	tests := []struct {
		name            string
		stockHistory    []models.Stock
		expectedContent []string // Strings that should be present in rationale
	}{
		{
			name:            "Empty history",
			stockHistory:    []models.Stock{},
			expectedContent: []string{"Sin datos históricos disponibles"},
		},
		{
			name: "Single stock with all fields",
			stockHistory: []models.Stock{
				testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Goldman Sachs", "upgrade", "hold", "buy", "150.00", "175.00", 0),
			},
			expectedContent: []string{"Goldman Sachs", "upgrade", "hold", "buy", "Precio objetivo"},
		},
		{
			name: "Stock with missing brokerage",
			stockHistory: []models.Stock{
				testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "", "upgrade", "hold", "buy", "150.00", "175.00", 0),
			},
			expectedContent: []string{"fuente no especificada", "upgrade", "hold", "buy"},
		},
		{
			name: "Stock with missing action",
			stockHistory: []models.Stock{
				testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Goldman Sachs", "", "hold", "buy", "150.00", "175.00", 0),
			},
			expectedContent: []string{"Goldman Sachs", "acción no especificada", "hold", "buy"},
		},
		{
			name: "Stock with missing ratings",
			stockHistory: []models.Stock{
				testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Goldman Sachs", "upgrade", "", "", "150.00", "175.00", 0),
			},
			expectedContent: []string{"Goldman Sachs", "upgrade", "sin rating previo", "rating no especificado"},
		},
		{
			name: "Stock with invalid target prices",
			stockHistory: []models.Stock{
				testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Goldman Sachs", "upgrade", "hold", "buy", "invalid", "also-invalid", 0),
			},
			expectedContent: []string{"Goldman Sachs", "upgrade", "hold", "buy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rationale := analyzer.generateRationale(tt.stockHistory)

			assert.NotEmpty(t, rationale, "Rationale should not be empty")

			for _, content := range tt.expectedContent {
				assert.Contains(t, rationale, content, "Rationale should contain: %s", content)
			}
		})
	}
}

func TestStockAnalyzer_groupStocksByTicker(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	stocks := []models.Stock{
		testhelpers.CreateTestStock("AAPL", "Apple Inc", 0),
		testhelpers.CreateTestStock("GOOGL", "Alphabet Inc", -time.Hour),
		testhelpers.CreateTestStock("AAPL", "Apple Inc", -2*time.Hour), // Another AAPL
		testhelpers.CreateTestStock("MSFT", "Microsoft Corp", -3*time.Hour),
	}

	stockMap := analyzer.groupStocksByTicker(stocks)

	assert.Len(t, stockMap, 3, "Should group into 3 tickers")
	assert.Len(t, stockMap["AAPL"], 2, "AAPL should have 2 entries")
	assert.Len(t, stockMap["GOOGL"], 1, "GOOGL should have 1 entry")
	assert.Len(t, stockMap["MSFT"], 1, "MSFT should have 1 entry")

	// Verify all original stocks are preserved
	totalStocks := 0
	for _, stocks := range stockMap {
		totalStocks += len(stocks)
	}
	assert.Equal(t, 4, totalStocks, "All stocks should be preserved")
}

func TestStockAnalyzer_groupStocksByTicker_EmptyInput(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	stocks := []models.Stock{}
	stockMap := analyzer.groupStocksByTicker(stocks)

	assert.Empty(t, stockMap, "Should return empty map for empty input")
}

// Integration-style test combining multiple components
func TestStockAnalyzer_Integration_FullWorkflow(t *testing.T) {
	mockRepo := mocks.NewMockStockRepository()
	analyzer := NewStockAnalyzer(mockRepo)

	// Create realistic test scenario
	baseTime := time.Now()
	testStocks := []models.Stock{
		// AAPL upgrade trend
		testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Goldman Sachs", "upgrade", "hold", "buy", "150.00", "175.00", 0),
		testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "Morgan Stanley", "upgrade", "neutral", "hold", "140.00", "150.00", -24*time.Hour),
		testhelpers.CreateTestStockFull("AAPL", "Apple Inc", "JP Morgan", "upgrade", "sell", "neutral", "120.00", "140.00", -48*time.Hour),

		// GOOGL single strong recommendation
		testhelpers.CreateTestStockFull("GOOGL", "Alphabet Inc", "Goldman Sachs", "upgrade", "hold", "strong buy", "2500.00", "2800.00", -time.Hour),

		// TSLA downgrade trend
		testhelpers.CreateTestStockFull("TSLA", "Tesla Inc", "Morgan Stanley", "downgrade", "buy", "hold", "800.00", "750.00", -2*time.Hour),
		testhelpers.CreateTestStockFull("TSLA", "Tesla Inc", "Goldman Sachs", "downgrade", "hold", "sell", "750.00", "600.00", -3*time.Hour),
	}

	// Set specific times
	for i := range testStocks {
		testStocks[i].Time = baseTime.Add(time.Duration(-i) * time.Hour)
	}

	// Setup mock to return testStocks for the filter
	mockRepo.On("GetFilteredStocks", mock.AnythingOfType("models.RecommendationFilter")).Return(testStocks, nil)

	// Test the full workflow
	filter := models.RecommendationFilter{
		Limit: 10,
	}

	recommendations, err := analyzer.GetFilteredRecommendations(filter)
	require.NoError(t, err)
	assert.Len(t, recommendations, 3, "Should have recommendations for 3 tickers")

	// Verify recommendations are properly sorted by score
	for i := 1; i < len(recommendations); i++ {
		assert.GreaterOrEqual(t, recommendations[i-1].Score, recommendations[i].Score,
			"Recommendations should be sorted by score")
	}

	// Verify that all scores are reasonable
	for _, rec := range recommendations {
		assert.GreaterOrEqual(t, rec.Score, 0.0, "Score should be non-negative")
		assert.LessOrEqual(t, rec.Score, 100.0, "Score should not exceed 100")
		assert.NotEmpty(t, rec.AnalysisRationale, "Should have analysis rationale")
	}

	// Verify that recommendations have expected latest ratings
	var aaplFound, googlFound, tslaFound bool
	for _, rec := range recommendations {
		switch rec.Ticker {
		case "AAPL":
			aaplFound = true
			assert.Equal(t, "buy", rec.LatestRating, "AAPL latest rating should be buy")
		case "GOOGL":
			googlFound = true
			assert.Equal(t, "strong buy", rec.LatestRating, "GOOGL latest rating should be strong buy")
		case "TSLA":
			tslaFound = true
			assert.Equal(t, "hold", rec.LatestRating, "TSLA latest rating should be hold")
		}
	}

	assert.True(t, aaplFound, "AAPL should be in recommendations")
	assert.True(t, googlFound, "GOOGL should be in recommendations")
	assert.True(t, tslaFound, "TSLA should be in recommendations")
}
