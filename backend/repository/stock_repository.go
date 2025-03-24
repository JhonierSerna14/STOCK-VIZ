package repository

import (
	"strings"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) SaveStocks(stocks []models.Stock) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "ticker"},
			{Name: "time"},
		},
		DoNothing: true,
	}).CreateInBatches(stocks, 100).Error
}

func (r *StockRepository) GetAllStocks() ([]models.Stock, error) {
	var stocks []models.Stock
	result := r.db.Find(&stocks)
	return stocks, result.Error
}

func (r *StockRepository) DeleteAllStocks() error {
	return r.db.Exec("DELETE FROM stocks").Error
}

// GetStocksPaginated recupera stocks con paginación
func (r *StockRepository) GetStocksPaginated(page, limit int, query string) ([]models.Stock, int64, error) {
	var stocks []models.Stock
	var total int64

	// Calcular el offset basado en la página y el límite
	offset := (page - 1) * limit

	// Base de la consulta
	dbQuery := r.db.Model(&models.Stock{})

	if query != "" {
		dbQuery = dbQuery.Where("LOWER(ticker) LIKE ? OR LOWER(company) LIKE ?", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%")
	}

	// Obtener total de registros con el filtro aplicado
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener registros paginados con el filtro aplicado
	result := dbQuery.Limit(limit).Offset(offset).Find(&stocks)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return stocks, total, nil
}

// GetFilteredStocks recupera stocks filtrados según los criterios especificados
func (r *StockRepository) GetFilteredStocks(filter models.RecommendationFilter) ([]models.Stock, error) {
	query := r.db.Model(&models.Stock{})

	// Aplicar filtro por ticker
	if filter.Ticker != "" {
		query = query.Where("ticker = ?", filter.Ticker)
	}

	// Aplicar filtro por fechas
	dateFrom, dateTo, err := filter.ParseDates()
	if err != nil {
		return nil, err
	}

	if dateFrom != nil {
		query = query.Where("time >= ?", dateFrom)
	}
	if dateTo != nil {
		query = query.Where("time <= ?", dateTo)
	}

	// Aplicar filtro por rating
	if filter.Rating != "" {
		query = query.Where("rating_to = ?", filter.Rating)
	}

	// Obtener los stocks filtrados
	var stocks []models.Stock
	if err := query.Find(&stocks).Error; err != nil {
		return nil, err
	}
	return stocks, nil
}
