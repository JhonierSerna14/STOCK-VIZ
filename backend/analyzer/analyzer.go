// Contiene la estructura principal y sus métodos básicos
package analyzer

import (
	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/calculator"
	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/formatter"
	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/scoring"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

type StockAnalyzer struct {
	repository StockRepository
	calculator *calculator.ScoreCalculator
	formatter  *formatter.NumberFormatter
}

type StockRepository interface {
	GetAllStocks() ([]models.Stock, error)
}

func NewStockAnalyzer(repo StockRepository) *StockAnalyzer {
	return &StockAnalyzer{
		repository: repo,
		calculator: calculator.NewScoreCalculator(scoring.DefaultWeights),
		formatter:  &formatter.NumberFormatter{},
	}
}
