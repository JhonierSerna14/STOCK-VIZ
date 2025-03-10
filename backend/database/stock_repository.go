package database

import (
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
		DoUpdates: clause.AssignmentColumns([]string{
			"company", "brokerage", "action",
			"rating_from", "rating_to",
			"target_from", "target_to",
		}),
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
