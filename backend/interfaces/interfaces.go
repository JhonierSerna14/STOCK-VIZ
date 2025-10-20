// Package interfaces defines the interfaces used throughout the application
package interfaces

import "github.com/JhonierSerna14/STOCK-VIZ/models"

// StockRepositoryInterface defines the contract for stock repository operations
type StockRepositoryInterface interface {
	SaveStocks(stocks []models.Stock) error
	GetAllStocks() ([]models.Stock, error)
	DeleteAllStocks() error
	GetStocksPaginated(page, limit int, query string) ([]models.Stock, int64, error)
	GetFilteredStocks(filter models.RecommendationFilter) ([]models.Stock, error)
}
