package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/mocks"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func createTestStock(ticker, company string) models.Stock {
	fixedTime := time.Date(2023, 10, 15, 12, 0, 0, 0, time.UTC)
	return models.Stock{
		Ticker:     ticker,
		Company:    company,
		Brokerage:  "Test Brokerage",
		Action:     "buy",
		RatingFrom: "neutral",
		RatingTo:   "overweight",
		TargetFrom: "100",
		TargetTo:   "150",
		Time:       fixedTime,
	}
}

func TestNewStockService(t *testing.T) {
	mockRepo := &mocks.MockStockRepository{}
	config := StockServiceConfig{
		Repository:   mockRepo,
		SyncInterval: time.Minute,
		APIToken:     "test-token",
		BaseURL:      "http://test.com",
	}

	service := NewStockService(config)
	require.NotNil(t, service)
	assert.Equal(t, config.Repository, service.config.Repository)
	assert.Equal(t, config.APIToken, service.config.APIToken)
	assert.Equal(t, config.BaseURL, service.config.BaseURL)
	assert.NotNil(t, service.analyzer)
}

func TestStockService_GetStocks(t *testing.T) {
	tests := []struct {
		name         string
		nextPage     string
		serverSetup  func() *httptest.Server
		setupMock    func(*mocks.MockStockRepository)
		wantError    bool
		expectedData *models.StockResponse
	}{
		{
			name:     "successful request without pagination",
			nextPage: "",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verificar headers
					assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
					assert.Equal(t, "GET", r.Method)

					stock1 := createTestStock("AAPL", "Apple Inc.")
					stock2 := createTestStock("GOOGL", "Alphabet Inc.")

					response := models.StockResponse{
						Items:    []models.Stock{stock1, stock2},
						NextPage: "",
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("SaveStocks", mock.AnythingOfType("[]models.Stock")).Return(nil)
			},
			wantError: false,
			expectedData: &models.StockResponse{
				Items:    []models.Stock{createTestStock("AAPL", "Apple Inc."), createTestStock("GOOGL", "Alphabet Inc.")},
				NextPage: "",
			},
		},
		{
			name:     "successful request with pagination",
			nextPage: "page2",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Verificar query parameter
					assert.Equal(t, "page2", r.URL.Query().Get("next_page"))

					stock1 := createTestStock("TSLA", "Tesla Inc.")

					response := models.StockResponse{
						Items:    []models.Stock{stock1},
						NextPage: "page3",
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("SaveStocks", mock.AnythingOfType("[]models.Stock")).Return(nil)
			},
			wantError: false,
			expectedData: &models.StockResponse{
				Items:    []models.Stock{createTestStock("TSLA", "Tesla Inc.")},
				NextPage: "page3",
			},
		},
		{
			name:     "server error",
			nextPage: "",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				// No se espera que se llame a SaveStocks
			},
			wantError:    true,
			expectedData: nil,
		},
		{
			name:     "invalid JSON response",
			nextPage: "",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("invalid json"))
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				// No se espera que se llame a SaveStocks
			},
			wantError:    true,
			expectedData: nil,
		},
		{
			name:     "repository save error",
			nextPage: "",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					stock1 := createTestStock("MSFT", "Microsoft Corp.")

					response := models.StockResponse{
						Items:    []models.Stock{stock1},
						NextPage: "",
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(response)
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("SaveStocks", mock.AnythingOfType("[]models.Stock")).Return(errors.New("database error"))
			},
			wantError:    true,
			expectedData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			server := tt.serverSetup()
			defer server.Close()

			mockRepo := &mocks.MockStockRepository{}
			tt.setupMock(mockRepo)

			config := StockServiceConfig{
				Repository:   mockRepo,
				SyncInterval: time.Minute,
				APIToken:     "test-token",
				BaseURL:      server.URL,
			}

			service := &StockService{
				config: config,
			}

			// Execute
			result, err := service.GetStocks(tt.nextPage)

			// Assert
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Items, len(tt.expectedData.Items))
				assert.Equal(t, tt.expectedData.NextPage, result.NextPage)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_GetAllStocks(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*mocks.MockStockRepository)
		wantError   bool
		expectedLen int
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *mocks.MockStockRepository) {
				stocks := []models.Stock{
					createTestStock("AAPL", "Apple Inc."),
					createTestStock("GOOGL", "Alphabet Inc."),
				}
				m.On("GetAllStocks").Return(stocks, nil)
			},
			wantError:   false,
			expectedLen: 2,
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("GetAllStocks").Return([]models.Stock{}, errors.New("database error"))
			},
			wantError:   true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockStockRepository{}
			tt.setupMock(mockRepo)

			service := &StockService{
				config: StockServiceConfig{Repository: mockRepo},
			}

			result, err := service.GetAllStocks()

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_GetAllStocksPaginated(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		limit         int
		query         string
		setupMock     func(*mocks.MockStockRepository)
		wantError     bool
		expectedLen   int
		expectedTotal int64
	}{
		{
			name:  "successful pagination",
			page:  1,
			limit: 10,
			query: "AAPL",
			setupMock: func(m *mocks.MockStockRepository) {
				stocks := []models.Stock{
					createTestStock("AAPL", "Apple Inc."),
				}
				m.On("GetStocksPaginated", 1, 10, "AAPL").Return(stocks, int64(1), nil)
			},
			wantError:     false,
			expectedLen:   1,
			expectedTotal: 1,
		},
		{
			name:  "repository error",
			page:  1,
			limit: 10,
			query: "",
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("GetStocksPaginated", 1, 10, "").Return([]models.Stock{}, int64(0), errors.New("database error"))
			},
			wantError:     true,
			expectedLen:   0,
			expectedTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockStockRepository{}
			tt.setupMock(mockRepo)

			service := &StockService{
				config: StockServiceConfig{Repository: mockRepo},
			}

			result, total, err := service.GetAllStocksPaginated(tt.page, tt.limit, tt.query)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
				assert.Equal(t, tt.expectedTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_DeleteAllStocks(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*mocks.MockStockRepository)
		wantError bool
	}{
		{
			name: "successful deletion",
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("DeleteAllStocks").Return(nil)
			},
			wantError: false,
		},
		{
			name: "repository error",
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("DeleteAllStocks").Return(errors.New("database error"))
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockStockRepository{}
			tt.setupMock(mockRepo)

			service := &StockService{
				config: StockServiceConfig{Repository: mockRepo},
			}

			err := service.DeleteAllStocks()

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_SyncStocksWithAPI(t *testing.T) {
	tests := []struct {
		name          string
		serverSetup   func() *httptest.Server
		setupMock     func(*mocks.MockStockRepository)
		wantError     bool
		expectedCount int
	}{
		{
			name: "successful sync single page",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					stock1 := createTestStock("AAPL", "Apple Inc.")
					stock2 := createTestStock("GOOGL", "Alphabet Inc.")

					response := models.StockResponse{
						Items:    []models.Stock{stock1, stock2},
						NextPage: "",
					}
					json.NewEncoder(w).Encode(response)
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("SaveStocks", mock.AnythingOfType("[]models.Stock")).Return(nil)
			},
			wantError:     false,
			expectedCount: 2,
		},
		{
			name: "successful sync multiple pages",
			serverSetup: func() *httptest.Server {
				callCount := 0
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					callCount++
					var response models.StockResponse

					if callCount == 1 {
						stock1 := createTestStock("AAPL", "Apple Inc.")
						response = models.StockResponse{
							Items:    []models.Stock{stock1},
							NextPage: "page2",
						}
					} else {
						stock2 := createTestStock("GOOGL", "Alphabet Inc.")
						response = models.StockResponse{
							Items:    []models.Stock{stock2},
							NextPage: "",
						}
					}
					json.NewEncoder(w).Encode(response)
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("SaveStocks", mock.AnythingOfType("[]models.Stock")).Return(nil).Twice()
			},
			wantError:     false,
			expectedCount: 2,
		},
		{
			name: "API error",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			setupMock: func(m *mocks.MockStockRepository) {
				// No se espera ninguna llamada
			},
			wantError:     true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverSetup()
			defer server.Close()

			mockRepo := &mocks.MockStockRepository{}
			tt.setupMock(mockRepo)

			config := StockServiceConfig{
				Repository:   mockRepo,
				SyncInterval: time.Minute,
				APIToken:     "test-token",
				BaseURL:      server.URL,
			}

			service := &StockService{
				config: config,
			}

			count, err := service.SyncStocksWithAPI()

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, count)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_GetRecommendations(t *testing.T) {
	tests := []struct {
		name        string
		filter      models.RecommendationFilter
		setupMock   func(*mocks.MockStockRepository)
		wantError   bool
		expectedLen int
	}{
		{
			name: "successful recommendations",
			filter: models.RecommendationFilter{
				Limit: 10,
			},
			setupMock: func(m *mocks.MockStockRepository) {
				stocks := []models.Stock{
					createTestStock("AAPL", "Apple Inc."),
					createTestStock("GOOGL", "Alphabet Inc."),
				}
				m.On("GetFilteredStocks", mock.AnythingOfType("models.RecommendationFilter")).Return(stocks, nil)
			},
			wantError:   false,
			expectedLen: 2,
		},
		{
			name: "repository error",
			filter: models.RecommendationFilter{
				Limit: 10,
			},
			setupMock: func(m *mocks.MockStockRepository) {
				m.On("GetFilteredStocks", mock.AnythingOfType("models.RecommendationFilter")).Return([]models.Stock{}, errors.New("database error"))
			},
			wantError:   true,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockStockRepository{}
			tt.setupMock(mockRepo)

			config := StockServiceConfig{
				Repository:   mockRepo,
				SyncInterval: time.Minute,
				APIToken:     "test-token",
				BaseURL:      "http://test.com",
			}

			service := NewStockService(config)

			result, err := service.GetRecommendations(tt.filter)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStockService_StartPeriodicSync(t *testing.T) {
	// Esta prueba verifica que StartPeriodicSync no cause panic
	// y que configure correctamente el servicio
	t.Run("start periodic sync", func(t *testing.T) {
		mockRepo := &mocks.MockStockRepository{}

		config := StockServiceConfig{
			Repository:   mockRepo,
			SyncInterval: 10 * time.Millisecond, // Intervalo muy pequeño para test
			APIToken:     "test-token",
			BaseURL:      "http://test.com",
		}

		service := &StockService{
			config: config,
		}

		// No debería causar panic
		assert.NotPanics(t, func() {
			service.StartPeriodicSync(config.SyncInterval)
		})

		// Dar tiempo para que arranque la goroutine
		time.Sleep(50 * time.Millisecond)
	})
}
