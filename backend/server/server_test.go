package server

import (
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/JhonierSerna14/STOCK-VIZ/api"
"github.com/JhonierSerna14/STOCK-VIZ/config"
"github.com/JhonierSerna14/STOCK-VIZ/repository"
"github.com/JhonierSerna14/STOCK-VIZ/service"
)

func createTestServer() *Server {
repo := repository.NewStockRepository(nil)
stockService := service.NewStockService(service.StockServiceConfig{
Repository:   repo,
SyncInterval: time.Minute,
APIToken:     "test-token",
BaseURL:      "http://test.com",
})

apiServer := api.NewAPIWithService(stockService)

cfg := &config.Config{
ServerPort: "8080",
CORSSettings: config.CORSConfig{
AllowedOrigins:   []string{"*"},
AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
AllowedHeaders:   []string{"*"},
AllowCredentials: true,
},
}

server := NewServer(apiServer, cfg)
server.SetupRoutes()
return server
}

func TestNewServer(t *testing.T) {
server := createTestServer()

if server == nil {
t.Error("Expected server to be non-nil")
}
if server.router == nil {
t.Error("Expected router to be non-nil")
}
if server.api == nil {
t.Error("Expected api to be non-nil")
}
if server.cors == nil {
t.Error("Expected cors to be non-nil")
}
if server.port != "8080" {
t.Errorf("Expected port 8080, got %s", server.port)
}
}

func TestServer_CORSHeaders(t *testing.T) {
server := createTestServer()

req, err := http.NewRequest("OPTIONS", "/api/stocks/all", nil)
if err != nil {
t.Fatalf("Failed to create request: %v", err)
}
req.Header.Set("Origin", "http://localhost:3000")
req.Header.Set("Access-Control-Request-Method", "GET")

rr := httptest.NewRecorder()
handler := server.cors.Handler(server.router)
handler.ServeHTTP(rr, req)

origin := rr.Header().Get("Access-Control-Allow-Origin")
if origin != "*" {
t.Errorf("Expected Access-Control-Allow-Origin *, got %s", origin)
}
}

func TestServer_CustomPort(t *testing.T) {
repo := repository.NewStockRepository(nil)
stockService := service.NewStockService(service.StockServiceConfig{
Repository:   repo,
SyncInterval: time.Minute,
APIToken:     "test-token",
BaseURL:      "http://test.com",
})

apiServer := api.NewAPIWithService(stockService)

cfg := &config.Config{
ServerPort: "9090",
CORSSettings: config.CORSConfig{
AllowedOrigins:   []string{"http://localhost:3000"},
AllowedMethods:   []string{"GET", "POST"},
AllowedHeaders:   []string{"Content-Type"},
AllowCredentials: false,
},
}

server := NewServer(apiServer, cfg)

if server.port != "9090" {
t.Errorf("Expected port 9090, got %s", server.port)
}
}

func TestServer_RouterNotNil(t *testing.T) {
server := createTestServer()
if server.router == nil {
t.Error("Router should not be nil after setup")
}
}
