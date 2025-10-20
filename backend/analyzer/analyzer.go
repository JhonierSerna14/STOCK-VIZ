// analyzer contiene la estructura principal y sus métodos básicos
package analyzer

import (
	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/calculator"
	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/formatter"
	"github.com/JhonierSerna14/STOCK-VIZ/analyzer/scoring"
	"github.com/JhonierSerna14/STOCK-VIZ/interfaces"
)

type StockAnalyzer struct {
	repository interfaces.StockRepositoryInterface
	calculator *calculator.ScoreCalculator
	formatter  *formatter.NumberFormatter
}

func NewStockAnalyzer(repo interfaces.StockRepositoryInterface) *StockAnalyzer {
	return &StockAnalyzer{
		repository: repo,
		calculator: calculator.NewScoreCalculator(scoring.DefaultWeights),
		formatter:  &formatter.NumberFormatter{},
	}
}
