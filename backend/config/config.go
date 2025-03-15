// Package config proporciona funcionalidades para gestionar la configuración de la aplicación.
package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config almacena la configuración de la aplicación
type Config struct {
	DatabaseURL  string
	SyncInterval time.Duration
	ServerPort   string
	CORSSettings CORSConfig
}

// CORSConfig almacena la configuración de CORS
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() (*Config, error) {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error cargando archivo .env: %v\n", err)
		// Continuamos incluso si no hay archivo .env, ya que las variables pueden estar configuradas directamente en el entorno
	}

	// Obtener la URL de la base de datos
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	// Obtener el puerto del servidor (valor predeterminado: 8080)
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// Obtener el intervalo de sincronización
	syncInterval := getEnvSyncInterval()

	// Configurar CORS
	corsConfig := CORSConfig{
		AllowedOrigins:   []string{"*"}, // Permite cualquier origen
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Permite cualquier header
		AllowCredentials: true,
	}

	return &Config{
		DatabaseURL:  databaseURL,
		SyncInterval: syncInterval,
		ServerPort:   serverPort,
		CORSSettings: corsConfig,
	}, nil
}

// getEnvSyncInterval obtiene el intervalo de sincronización desde las variables de entorno
func getEnvSyncInterval() time.Duration {
	syncIntervalStr := os.Getenv("SYNC_INTERVAL_MINUTES")
	var syncInterval time.Duration
	if syncIntervalStr != "" {
		minutes, err := strconv.Atoi(syncIntervalStr)
		if err != nil || minutes <= 0 {
			log.Printf("Valor inválido para SYNC_INTERVAL_MINUTES: %s. Usando valor predeterminado de 60 minutos.\n", syncIntervalStr)
			syncInterval = 60 * time.Minute
		} else {
			syncInterval = time.Duration(minutes) * time.Minute
		}
	} else {
		// Valor predeterminado: sincronizar cada hora
		syncInterval = 60 * time.Minute
	}
	return syncInterval
}
