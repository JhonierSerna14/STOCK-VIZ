package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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

// Método para crear un servicio con sincronización automática
func NewStockServiceWithSync(repo *database.StockRepository, syncInterval time.Duration) *StockService {
	service := NewStockService(repo)

	// Iniciar la sincronización periódica
	if syncInterval > 0 {
		service.StartPeriodicSync(syncInterval)
	}

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

func (s *StockService) GetRecommendations(filter models.RecommendationFilter) ([]models.StockRecommendation, error) {
	// Si no se proporciona límite, usar valor por defecto
	limit := filter.Limit
	if limit <= 0 {
		limit = 5
	}

	return s.analyzer.GetFilteredRecommendations(filter, limit)
}

// MigrateAllStocks recupera todos los stocks de la API externa
// y los guarda en la base de datos local de forma recursiva,
// hasta que ya no haya más páginas disponibles.
func (s *StockService) MigrateAllStocks() (int, error) {
	var nextPage string
	var totalItems int

	for {
		stockResponse, err := s.GetStocks(nextPage)
		if err != nil {
			return totalItems, fmt.Errorf("error obteniendo stocks: %v", err)
		}

		// Contamos cuántos items hemos procesado
		totalItems += len(stockResponse.Items)

		// Si no hay más páginas, terminamos el proceso
		if stockResponse.NextPage == "" {
			break
		}

		// Preparamos la siguiente página
		nextPage = stockResponse.NextPage
	}

	return totalItems, nil
}

// SyncStocksWithAPI sincroniza los stocks de la API externa con la base de datos local.
// Devuelve el número de nuevos stocks añadidos.
func (s *StockService) SyncStocksWithAPI() (int, error) {
	var nextPage string
	var totalNewItems int

	for {
		stockResponse, err := s.GetStocks(nextPage)
		if err != nil {
			return totalNewItems, fmt.Errorf("error obteniendo stocks: %v", err)
		}

		// Contamos cuántos items hemos procesado
		totalNewItems += len(stockResponse.Items)

		// Si no hay más páginas, terminamos el proceso
		if stockResponse.NextPage == "" {
			break
		}

		// Preparamos la siguiente página
		nextPage = stockResponse.NextPage
	}

	return totalNewItems, nil
}

// StartPeriodicSync inicia un proceso en segundo plano que sincroniza
// periódicamente la base de datos local con la API externa.
func (s *StockService) StartPeriodicSync(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				count, err := s.SyncStocksWithAPI()
				if err != nil {
					fmt.Printf("Error en la sincronización periódica: %v\n", err)
				} else {
					fmt.Printf("Sincronización completada: %d stocks procesados\n", count)
				}
			}
		}
	}()
	fmt.Printf("Sincronización periódica iniciada con intervalo de %v\n", interval)
}

// GetAllStocksPaginated obtiene stocks paginados de la base de datos
func (s *StockService) GetAllStocksPaginated(page, limit int) ([]models.Stock, int64, error) {
	return s.repository.GetStocksPaginated(page, limit)
}
