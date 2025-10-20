// Package mocks provides mock implementations for testing
package mocks

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/interfaces"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/stretchr/testify/mock"
)

// Ensure MockStockRepository implements the interface
var _ interfaces.StockRepositoryInterface = (*MockStockRepository)(nil)

// MockStockRepository is a mock implementation of StockRepository for testing
type MockStockRepository struct {
	mock.Mock
}

// NewMockStockRepository creates a new mock repository
func NewMockStockRepository() *MockStockRepository {
	return &MockStockRepository{}
}

// SaveStocks mocks the SaveStocks method
func (m *MockStockRepository) SaveStocks(stocks []models.Stock) error {
	args := m.Called(stocks)
	return args.Error(0)
}

// GetAllStocks mocks the GetAllStocks method
func (m *MockStockRepository) GetAllStocks() ([]models.Stock, error) {
	args := m.Called()
	return args.Get(0).([]models.Stock), args.Error(1)
}

// DeleteAllStocks mocks the DeleteAllStocks method
func (m *MockStockRepository) DeleteAllStocks() error {
	args := m.Called()
	return args.Error(0)
}

// GetStocksPaginated mocks the GetStocksPaginated method
func (m *MockStockRepository) GetStocksPaginated(page, limit int, query string) ([]models.Stock, int64, error) {
	args := m.Called(page, limit, query)
	return args.Get(0).([]models.Stock), args.Get(1).(int64), args.Error(2)
}

// GetFilteredStocks mocks the GetFilteredStocks method
func (m *MockStockRepository) GetFilteredStocks(filter models.RecommendationFilter) ([]models.Stock, error) {
	args := m.Called(filter)
	return args.Get(0).([]models.Stock), args.Error(1)
}

// MockHTTPServer creates a mock HTTP server for testing API calls
type MockHTTPServer struct {
	*httptest.Server
	responses   map[string]string
	statusCodes map[string]int
}

// NewMockHTTPServer creates a new mock HTTP server
func NewMockHTTPServer() *MockHTTPServer {
	mock := &MockHTTPServer{
		responses:   make(map[string]string),
		statusCodes: make(map[string]int),
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if r.URL.RawQuery != "" {
			path += "?" + r.URL.RawQuery
		}

		// Check if we have a specific response for this path
		if response, exists := mock.responses[path]; exists {
			statusCode := mock.statusCodes[path]
			if statusCode == 0 {
				statusCode = http.StatusOK
			}
			w.WriteHeader(statusCode)
			w.Write([]byte(response))
			return
		}

		// Default response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"items":[],"next_page":null}`))
	}))

	return mock
}

// SetResponse sets a response for a specific path
func (m *MockHTTPServer) SetResponse(path, response string) {
	m.responses[path] = response
}

// SetStatusCode sets a status code for a specific path
func (m *MockHTTPServer) SetStatusCode(path string, statusCode int) {
	m.statusCodes[path] = statusCode
}

// SetStockResponse sets a stock response for a specific path
func (m *MockHTTPServer) SetStockResponse(path string, stocks []models.Stock, nextPage *string) {
	// This is a simplified JSON marshaling - in real tests you might want to use json.Marshal
	jsonResponse := `{"items":[`
	for i, stock := range stocks {
		if i > 0 {
			jsonResponse += ","
		}
		jsonResponse += fmt.Sprintf(`{
			"ticker":"%s",
			"company":"%s",
			"brokerage":"%s",
			"action":"%s",
			"rating_from":"%s",
			"rating_to":"%s",
			"target_from":"%s",
			"target_to":"%s",
			"time":"%s"
		}`, stock.Ticker, stock.Company, stock.Brokerage, stock.Action,
			stock.RatingFrom, stock.RatingTo, stock.TargetFrom, stock.TargetTo,
			stock.Time.Format(time.RFC3339))
	}
	jsonResponse += `]`

	if nextPage != nil {
		jsonResponse += fmt.Sprintf(`,"next_page":"%s"`, *nextPage)
	} else {
		jsonResponse += `,"next_page":null`
	}
	jsonResponse += `}`

	m.SetResponse(path, jsonResponse)
}

// Close closes the mock server
func (m *MockHTTPServer) Close() {
	m.Server.Close()
}
