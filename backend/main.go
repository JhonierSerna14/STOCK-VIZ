// Package main es el punto de entrada principal para la aplicación STOCK-VIZ.
package main

import (
	"log"

	"github.com/JhonierSerna14/STOCK-VIZ/api"
	"github.com/JhonierSerna14/STOCK-VIZ/config"
	"github.com/JhonierSerna14/STOCK-VIZ/database"
	"github.com/JhonierSerna14/STOCK-VIZ/server"
	"github.com/JhonierSerna14/STOCK-VIZ/service"
)

// main inicializa y ejecuta la aplicación STOCK-VIZ.
func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración: %v\n", err)
	}

	// Inicializar conexión a la base de datos
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error inicializando base de datos: %v\n", err)
	}

	// Inicializar el repositorio y servicio con sincronización automática
	stockRepo := database.NewStockRepository(db)
	stockService := service.NewStockServiceWithSync(stockRepo, cfg.SyncInterval)

	// Inicializar la API con el servicio configurado
	a := api.NewAPIWithService(stockService)

	// Crear y configurar el servidor HTTP
	s := server.NewServer(a, cfg)
	s.SetupRoutes()

	// Iniciar el servidor
	if err := s.Start(); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v\n", err)
	}
}
