package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/JhonierSerna14/STOCK-VIZ/analyzer"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/JhonierSerna14/STOCK-VIZ/repository"
)

// StockService representa el servicio principal para gestionar stocks
type StockService struct {
	config   StockServiceConfig
	analyzer *analyzer.StockAnalyzer
}

// StockServiceConfig contiene la configuración necesaria para crear un StockService
type StockServiceConfig struct {
	Repository   *repository.StockRepository
	SyncInterval time.Duration
	APIToken     string
	BaseURL      string
}

// NewStockService crea una nueva instancia de StockService con la configuración proporcionada
func NewStockService(cfg StockServiceConfig) *StockService {
	service := &StockService{
		config:   cfg,
		analyzer: analyzer.NewStockAnalyzer(cfg.Repository),
	}

	service.StartPeriodicSync(cfg.SyncInterval)

	return service
}

// GetStocks obtiene stocks de la API externa
func (s *StockService) GetStocks(nextPage string) (*models.StockResponse, error) {
	url := s.config.BaseURL
	if nextPage != "" {
		url += "?next_page=" + nextPage
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+s.config.APIToken)

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
	if err := s.config.Repository.SaveStocks(stockResponse.Items); err != nil {
		return nil, fmt.Errorf("error guardando stocks: %v", err)
	}

	return &stockResponse, nil
}

// GetAllStocks obtiene todos los stocks de la base de datos
func (s *StockService) GetAllStocks() ([]models.Stock, error) {
	return s.config.Repository.GetAllStocks()
}

// GetAllStocksPaginated obtiene stocks paginados
func (s *StockService) GetAllStocksPaginated(page, limit int) ([]models.Stock, int64, error) {
	return s.config.Repository.GetStocksPaginated(page, limit)
}

// DeleteAllStocks elimina todos los stocks de la base de datos
func (s *StockService) DeleteAllStocks() error {
	return s.config.Repository.DeleteAllStocks()
}

// SyncStocksWithAPI sincroniza los stocks con la API externa
func (s *StockService) SyncStocksWithAPI() (int, error) {
	var nextPage string
	var totalNewItems int

	for {
		stockResponse, err := s.GetStocks(nextPage)
		if err != nil {
			return totalNewItems, fmt.Errorf("error obteniendo stocks: %v", err)
		}

		totalNewItems += len(stockResponse.Items)

		if stockResponse.NextPage == "" {
			break
		}

		nextPage = stockResponse.NextPage
	}

	return totalNewItems, nil
}

// StartPeriodicSync inicia la sincronización periódica
func (s *StockService) StartPeriodicSync(interval time.Duration) {
	// Iniciar sincronizaciones en segundo plano
	go func() {
		// Realizar la sincronización inicial inmediatamente
		count, err := s.SyncStocksWithAPI()
		if err != nil {
			fmt.Printf("Error en la sincronización inicial: %v\n", err)
		} else {
			fmt.Printf("Sincronización inicial completada: %d stocks procesados\n", count)
		}

		// Configurar el ticker para sincronizaciones periódicas
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				count, err := s.SyncStocksWithAPI()
				if err != nil {
					fmt.Printf("Error en la sincronización periódica: %v\n", err)
				} else {
					fmt.Printf("Sincronización periódica completada: %d stocks procesados\n", count)
				}
			}
		}
	}()

	fmt.Printf("Sincronización iniciada: primera sincronización en proceso, siguientes cada %v\n", interval)
}

// GetRecommendations obtiene recomendaciones filtradas
func (s *StockService) GetRecommendations(filter models.RecommendationFilter) ([]models.StockRecommendation, error) {
	return s.analyzer.GetFilteredRecommendations(filter)
}
