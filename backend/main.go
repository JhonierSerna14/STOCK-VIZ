package main

import (
	"log"
	"net/http"
	"os"

	"github.com/JhonierSerna14/STOCK-VIZ/api"
	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando archivo .env: %v\n", err)
	}

	// Inicializar conexión a la base de datos con GORM
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v\n", err)
	}

	// Agregar migración automática
	if err := db.AutoMigrate(&models.Stock{}); err != nil {
		log.Fatalf("Error en la migración de la base de datos: %v\n", err)
	}

	r := mux.NewRouter()

	// Usar db (GORM)
	a := api.NewAPI(db)

	// Register the routes
	a.RegisterRoutes(r)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
