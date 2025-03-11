package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/JhonierSerna14/STOCK-VIZ/analyzer"
	"github.com/JhonierSerna14/STOCK-VIZ/database"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
)

type StockService struct {
	baseURL    string
	token      string
	repository *database.StockRepository
	analyzer   *analyzer.StockAnalyzer
}

func NewStockService(repo *database.StockRepository) *StockService {
	token := os.Getenv("STOCK_API_TOKEN")
	if token == "" {
		fmt.Println("Advertencia: STOCK_API_TOKEN no está configurado")
	}

	baseURL := os.Getenv("STOCK_API_BASE_URL")
	if baseURL == "" {
		fmt.Println("Advertencia: STOCK_API_BASE_URL no está configurado")
	}

	service := &StockService{
		baseURL:    baseURL,
		token:      token,
		repository: repo,
	}

	service.analyzer = analyzer.NewStockAnalyzer(repo)
	return service
}

func (s *StockService) GetStocks(nextPage string) (*models.StockResponse, error) {
	url := s.baseURL
	if nextPage != "" {
		url += "?next_page=" + nextPage
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stockResponse models.StockResponse
	if err := json.Unmarshal(body, &stockResponse); err != nil {
		return nil, err
	}

	// Guardar los datos en la base de datos
	if err := s.repository.SaveStocks(stockResponse.Items); err != nil {
		return nil, fmt.Errorf("error guardando stocks: %v", err)
	}

	return &stockResponse, nil
}

func (s *StockService) GetAllStocks() ([]models.Stock, error) {
	return s.repository.GetAllStocks()
}

func (s *StockService) DeleteAllStocks() error {
	return s.repository.DeleteAllStocks()
}

func (s *StockService) GetRecommendations() ([]models.StockRecommendation, error) {
	return s.analyzer.GetTopRecommendations(5) // Retorna top 5 recomendaciones
}
