// Package server proporciona funcionalidades para configurar y ejecutar el servidor HTTP.
package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JhonierSerna14/STOCK-VIZ/api"
	"github.com/JhonierSerna14/STOCK-VIZ/config"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Server representa el servidor HTTP
type Server struct {
	router *mux.Router
	api    *api.API
	port   string
	cors   *cors.Cors
}

// NewServer crea una nueva instancia del servidor
func NewServer(api *api.API, cfg *config.Config) *Server {
	// Crear un nuevo router usando Gorilla Mux
	r := mux.NewRouter()

	// Configurar CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORSSettings.AllowedOrigins,
		AllowedMethods:   cfg.CORSSettings.AllowedMethods,
		AllowedHeaders:   cfg.CORSSettings.AllowedHeaders,
		AllowCredentials: cfg.CORSSettings.AllowCredentials,
	})

	return &Server{
		router: r,
		api:    api,
		port:   cfg.ServerPort,
		cors:   c,
	}
}

// SetupRoutes configura las rutas del servidor
func (s *Server) SetupRoutes() {
	// Registrar las rutas HTTP para la API
	s.api.RegisterRoutes(s.router)
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	// Envolver el router con el middleware CORS
	handler := s.cors.Handler(s.router)

	// Iniciar el servidor HTTP en el puerto configurado
	log.Printf("Iniciando servidor en el puerto :%s\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), handler)
}
