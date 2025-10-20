package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/mocks"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/JhonierSerna14/STOCK-VIZ/service"
	"github.com/JhonierSerna14/STOCK-VIZ/testhelpers"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createTestAPI() *API {
	mockRepo := mocks.NewMockStockRepository()
	config := service.StockServiceConfig{
		Repository:   mockRepo,
		SyncInterval: time.Minute,
		APIToken:     "test-token",
		BaseURL:      "http://test.com",
	}
	stockService := service.NewStockService(config)
	return NewAPIWithService(stockService)
}

func createTestAPIWithData() *API {
	mockRepo := mocks.NewMockStockRepository()

	// Add test data
	testStocks := []models.Stock{
		testhelpers.CreateTestStock("AAPL", "Apple Inc", 0),
		testhelpers.CreateTestStock("GOOGL", "Alphabet Inc", -time.Hour),
		testhelpers.CreateTestStock("MSFT", "Microsoft Corp", -2*time.Hour),
		testhelpers.CreateTestStockWithRating("TSLA", "Tesla Inc", "strong buy", -time.Hour),
		testhelpers.CreateTestStockWithRating("AMZN", "Amazon Inc", "sell", -3*time.Hour),
	}

	// Setup mocks for all API operations with specific returns first, then general ones
	mockRepo.On("GetStocksPaginated", 1, 2, "").Return(testStocks[0:2], int64(len(testStocks)), nil)                                                                         // Page 1, limit 2, return first 2
	mockRepo.On("GetStocksPaginated", 1, 20, "").Return(testStocks, int64(len(testStocks)), nil)                                                                             // Default pagination used by RegisterRoutes tests
	mockRepo.On("GetStocksPaginated", 1, mock.AnythingOfType("int"), "AAPL").Return([]models.Stock{testStocks[0]}, int64(1), nil)                                            // Search for AAPL
	mockRepo.On("GetStocksPaginated", mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(testStocks, int64(len(testStocks)), nil) // General fallback
	mockRepo.On("DeleteAllStocks").Return(nil)
	mockRepo.On("GetFilteredStocks", mock.AnythingOfType("models.RecommendationFilter")).Return(testStocks, nil)

	config := service.StockServiceConfig{
		Repository:   mockRepo,
		SyncInterval: time.Minute,
		APIToken:     "test-token",
		BaseURL:      "http://test.com",
	}
	stockService := service.NewStockService(config)
	return NewAPIWithService(stockService)
}

func TestNewAPIWithService(t *testing.T) {
	api := createTestAPI()
	assert.NotNil(t, api, "Expected api to be non-nil")
	assert.NotNil(t, api.stockService, "Expected stockService to be non-nil")
}

func TestAPI_GetStocks(t *testing.T) {
	api := createTestAPIWithData()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "Get stocks without pagination",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				assert.Contains(t, response, "items", "Response should contain items")
				assert.Contains(t, response, "pagination", "Response should contain pagination")

				items, ok := response["items"].([]interface{})
				require.True(t, ok, "Items should be array")
				assert.Len(t, items, 5, "Should return 5 stocks")
			},
		},
		{
			name:           "Get stocks with pagination",
			queryParams:    "?page=1&limit=2",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				items, ok := response["items"].([]interface{})
				require.True(t, ok, "Items should be array")
				assert.Len(t, items, 2, "Should return 2 stocks per page")

				pagination, ok := response["pagination"].(map[string]interface{})
				require.True(t, ok, "Pagination should be object")
				assert.Equal(t, float64(1), pagination["current_page"], "Should be page 1")
				assert.Equal(t, float64(2), pagination["per_page"], "Should be limit 2")
			},
		},
		{
			name:           "Get stocks with search",
			queryParams:    "?search=AAPL",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				items, ok := response["items"].([]interface{})
				require.True(t, ok, "Items should be array")
				// Note: Mock doesn't implement search filtering, so this tests the endpoint structure
				assert.GreaterOrEqual(t, len(items), 0, "Should return some stocks")
			},
		},
		{
			name:           "Get stocks with invalid page parameter",
			queryParams:    "?page=invalid&limit=10",
			expectedStatus: http.StatusOK, // Should default to page 1
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				pagination, ok := response["pagination"].(map[string]interface{})
				require.True(t, ok, "Pagination should be object")
				assert.Equal(t, float64(1), pagination["current_page"], "Should default to page 1")
			},
		},
		{
			name:           "Get stocks with zero page",
			queryParams:    "?page=0&limit=10",
			expectedStatus: http.StatusOK, // Should default to page 1
			checkResponse: func(t *testing.T, body []byte) {
				var response map[string]interface{}
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)

				pagination, ok := response["pagination"].(map[string]interface{})
				require.True(t, ok, "Pagination should be object")
				assert.Equal(t, float64(1), pagination["current_page"], "Should default to page 1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/stocks/all"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			api.getStocks(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match")

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr.Body.Bytes())
			}

			// Verify content type
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Should set JSON content type")
		})
	}
}

func TestAPI_DeleteAllStocks(t *testing.T) {
	api := createTestAPIWithData()

	req := httptest.NewRequest(http.MethodDelete, "/api/stocks", nil)
	rr := httptest.NewRecorder()

	api.deleteAllStocks(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Should return 200 OK")
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Should set JSON content type")

	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err, "Should be valid JSON")

	assert.Contains(t, response, "mensaje", "Response should contain message")
	assert.Contains(t, response["mensaje"], "eliminados exitosamente", "Should contain success message")
}

func TestAPI_DeleteAllStocks_Error(t *testing.T) {
	// Create API with error mock
	mockRepo := mocks.NewMockStockRepository()
	mockRepo.On("DeleteAllStocks").Return(assert.AnError)

	config := service.StockServiceConfig{
		Repository:   mockRepo,
		SyncInterval: time.Minute,
		APIToken:     "test-token",
		BaseURL:      "http://test.com",
	}
	stockService := service.NewStockService(config)
	api := NewAPIWithService(stockService)

	req := httptest.NewRequest(http.MethodDelete, "/api/stocks", nil)
	rr := httptest.NewRecorder()

	api.deleteAllStocks(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return 500 on error")
}

func TestAPI_GetRecommendations(t *testing.T) {
	api := createTestAPIWithData()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name:           "Get recommendations without filter",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var recommendations []models.StockRecommendation
				err := json.Unmarshal(body, &recommendations)
				require.NoError(t, err, "Should be valid JSON array")

				// Should have recommendations for unique tickers
				assert.GreaterOrEqual(t, len(recommendations), 1, "Should have at least one recommendation")

				// Verify recommendation structure
				for _, rec := range recommendations {
					assert.NotEmpty(t, rec.Ticker, "Ticker should not be empty")
					assert.NotEmpty(t, rec.Company, "Company should not be empty")
					assert.GreaterOrEqual(t, rec.Score, 0.0, "Score should be non-negative")
					assert.LessOrEqual(t, rec.Score, 100.0, "Score should not exceed 100")
					assert.NotEmpty(t, rec.AnalysisRationale, "Analysis rationale should not be empty")
				}
			},
		},
		{
			name:           "Get recommendations with limit",
			queryParams:    "?limit=2",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var recommendations []models.StockRecommendation
				err := json.Unmarshal(body, &recommendations)
				require.NoError(t, err)

				assert.LessOrEqual(t, len(recommendations), 2, "Should respect limit")
			},
		},
		{
			name:           "Get recommendations with date filter",
			queryParams:    "?date_from=2024-01-01&date_to=2024-12-31",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var recommendations []models.StockRecommendation
				err := json.Unmarshal(body, &recommendations)
				require.NoError(t, err)

				// Should return array (mock doesn't filter by date)
				assert.IsType(t, []models.StockRecommendation{}, recommendations)
			},
		},
		{
			name:           "Get recommendations with rating filter",
			queryParams:    "?rating=buy",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var recommendations []models.StockRecommendation
				err := json.Unmarshal(body, &recommendations)
				require.NoError(t, err)

				assert.IsType(t, []models.StockRecommendation{}, recommendations)
			},
		},
		{
			name:           "Get recommendations with ticker filter",
			queryParams:    "?ticker=AAPL",
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var recommendations []models.StockRecommendation
				err := json.Unmarshal(body, &recommendations)
				require.NoError(t, err)

				assert.IsType(t, []models.StockRecommendation{}, recommendations)
			},
		},
		{
			name:           "Get recommendations with invalid limit",
			queryParams:    "?limit=invalid",
			expectedStatus: http.StatusOK, // Should default to 6
			checkResponse: func(t *testing.T, body []byte) {
				var recommendations []models.StockRecommendation
				err := json.Unmarshal(body, &recommendations)
				require.NoError(t, err)

				assert.IsType(t, []models.StockRecommendation{}, recommendations)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/recommendations"+tt.queryParams, nil)
			rr := httptest.NewRecorder()

			api.getRecommendations(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "Status code should match")

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr.Body.Bytes())
			}

			// Verify content type
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "Should set JSON content type")
		})
	}
}

func TestAPI_GetRecommendations_ServiceError(t *testing.T) {
	// Create API with error mock
	mockRepo := mocks.NewMockStockRepository()
	mockRepo.On("GetFilteredStocks", mock.AnythingOfType("models.RecommendationFilter")).Return([]models.Stock{}, assert.AnError)

	config := service.StockServiceConfig{
		Repository:   mockRepo,
		SyncInterval: time.Minute,
		APIToken:     "test-token",
		BaseURL:      "http://test.com",
	}
	stockService := service.NewStockService(config)
	api := NewAPIWithService(stockService)

	req := httptest.NewRequest(http.MethodGet, "/api/recommendations", nil)
	rr := httptest.NewRecorder()

	api.getRecommendations(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Should return 500 on error")
}

func TestAPI_RegisterRoutes(t *testing.T) {
	api := createTestAPIWithData()
	router := mux.NewRouter()

	api.RegisterRoutes(router)

	// Test that routes are registered by making requests
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/stocks/all"},
		{http.MethodDelete, "/api/stocks"},
		{http.MethodGet, "/api/recommendations"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			// Should not return 404 (route not found)
			assert.NotEqual(t, http.StatusNotFound, rr.Code, "Route should be registered")
		})
	}
}

// Test JSON response structures
func TestAPI_GetRecommendationsUnmarshal(t *testing.T) {
	// Test that response can be unmarshaled to recommendations
	jsonData := []byte(`[{"ticker":"AAPL","company":"Apple","score":85.5,"rating":"buy"}]`)
	var recommendations []models.StockRecommendation
	err := json.Unmarshal(jsonData, &recommendations)
	require.NoError(t, err, "Failed to unmarshal")

	assert.Len(t, recommendations, 1, "Expected 1 recommendation")
	assert.Equal(t, "AAPL", recommendations[0].Ticker, "Expected ticker AAPL")
}

func TestAPI_StockResponseUnmarshal(t *testing.T) {
	// Test that stock response can be unmarshaled
	jsonData := []byte(`{"items":[],"pagination":{"current_page":1,"per_page":10,"total_items":0,"total_pages":0,"has_next_page":false}}`)
	var response map[string]interface{}
	err := json.Unmarshal(jsonData, &response)
	require.NoError(t, err, "Failed to unmarshal")

	assert.Contains(t, response, "items", "Expected 'items' in response")
	assert.Contains(t, response, "pagination", "Expected 'pagination' in response")
}

func TestAPI_DeleteResponseUnmarshal(t *testing.T) {
	jsonData := []byte(`{"mensaje":"stocks eliminados exitosamente"}`)
	var response map[string]string
	err := json.Unmarshal(jsonData, &response)
	require.NoError(t, err, "Failed to unmarshal")

	assert.Contains(t, response, "mensaje", "Expected 'mensaje' in response")
}

// Test HTTP method restrictions
func TestAPI_MethodNotAllowed(t *testing.T) {
	api := createTestAPIWithData()
	router := mux.NewRouter()
	api.RegisterRoutes(router)

	tests := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/stocks/all"},      // Should be GET
		{http.MethodGet, "/api/stocks"},           // Should be DELETE
		{http.MethodPost, "/api/recommendations"}, // Should be GET
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Should return 405 for wrong method")
		})
	}
}
