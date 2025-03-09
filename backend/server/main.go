package main

import (
	"log"
	"net/http"

	"github.com/JhonierSerna14/STOCK-VIZ/api"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Printf("Error cargando archivo .env: %v\n", err)
	}

	r := mux.NewRouter()

	// Create a new API instance using the constructor
	a := api.NewAPI()

	// Register the routes
	a.RegisterRoutes(r)

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
