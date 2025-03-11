// Package main es el punto de entrada principal para la aplicación STOCK-VIZ.
// Implementa la configuración de la base de datos, rutas de API y servidor HTTP.
package main

import (
	// Importaciones estándar
	"log"
	"net/http"
	"os"

	// Importaciones del proyecto
	"github.com/JhonierSerna14/STOCK-VIZ/api"
	"github.com/JhonierSerna14/STOCK-VIZ/models"

	// Importaciones de terceros
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// main inicializa y ejecuta la aplicación STOCK-VIZ.
// Configura la conexión a la base de datos, crea el router,
// registra las rutas de API y arranca el servidor HTTP.
func main() {
	// Cargar variables de entorno desde el archivo .env
	// Esto permite configurar la aplicación sin cambiar el código fuente
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando archivo .env: %v\n", err)
	}

	// Inicializar conexión a la base de datos con GORM
	// Utilizamos la URL de la base de datos desde las variables de entorno
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Configurar el nivel de logging para mostrar solo errores
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatalf("Error conectando a la base de datos: %v\n", err)
	}

	// Migración automática del esquema de base de datos
	// Esto crea o actualiza la tabla Stock según el modelo definido
	if err := db.AutoMigrate(&models.Stock{}); err != nil {
		log.Fatalf("Error en la migración de la base de datos: %v\n", err)
	}

	// Crear un nuevo router usando Gorilla Mux
	r := mux.NewRouter()

	// Inicializar la API con la conexión a la base de datos
	a := api.NewAPI(db)

	// Registrar las rutas HTTP para la API
	a.RegisterRoutes(r)

	// Iniciar el servidor HTTP en el puerto 8080
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
