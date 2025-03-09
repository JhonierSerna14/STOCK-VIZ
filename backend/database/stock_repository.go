package database

import (
	"database/sql"

	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

type StockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) *StockRepository {
	return &StockRepository{db: db}
}

func (r *StockRepository) SaveStocks(stocks []models.Stock) error {
	query := `
		INSERT INTO stocks (ticker, company, brokerage, action, rating_from, rating_to, target_from, target_to, time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (ticker, time) DO NOTHING`

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, stock := range stocks {
		_, err = stmt.Exec(
			stock.Ticker,
			stock.Company,
			stock.Brokerage,
			stock.Action,
			stock.RatingFrom,
			stock.RatingTo,
			stock.TargetFrom,
			stock.TargetTo,
			stock.Time,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
