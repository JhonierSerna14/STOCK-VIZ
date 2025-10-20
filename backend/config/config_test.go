package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestLoadConfigSuccess(t *testing.T) {
	// Setup environment variables
	setupTestEnv(t)
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if config == nil {
		t.Fatal("Expected config, got nil")
	}

	if config.DatabaseURL != "postgres://test:test@localhost/testdb" {
		t.Errorf("Expected DatabaseURL postgres://test:test@localhost/testdb, got %s", config.DatabaseURL)
	}
	if config.StockAPIToken != "test-token" {
		t.Errorf("Expected StockAPIToken test-token, got %s", config.StockAPIToken)
	}
	if config.BaseURL != "https://api.test.com" {
		t.Errorf("Expected BaseURL https://api.test.com, got %s", config.BaseURL)
	}
	if config.ServerPort != "8080" {
		t.Errorf("Expected ServerPort 8080, got %s", config.ServerPort)
	}
	if config.SyncInterval != 60*time.Minute {
		t.Errorf("Expected SyncInterval 60m, got %v", config.SyncInterval)
	}

	// Test CORS settings
	if len(config.CORSSettings.AllowedOrigins) != 1 || config.CORSSettings.AllowedOrigins[0] != "*" {
		t.Errorf("Expected AllowedOrigins [*], got %v", config.CORSSettings.AllowedOrigins)
	}
	methodsContainGet := false
	for _, method := range config.CORSSettings.AllowedMethods {
		if method == "GET" {
			methodsContainGet = true
			break
		}
	}
	if !methodsContainGet {
		t.Error("Expected GET in AllowedMethods")
	}
	if !config.CORSSettings.AllowCredentials {
		t.Error("Expected AllowCredentials to be true")
	}
}

func TestLoadConfigCustomPort(t *testing.T) {
	setupTestEnv(t)
	os.Setenv("SERVER_PORT", "9090")
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.ServerPort != "9090" {
		t.Errorf("Expected ServerPort 9090, got %s", config.ServerPort)
	}
}

func TestLoadConfigCustomSyncInterval(t *testing.T) {
	setupTestEnv(t)
	os.Setenv("SYNC_INTERVAL_MINUTES", "30")
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.SyncInterval != 30*time.Minute {
		t.Errorf("Expected SyncInterval 30m, got %v", config.SyncInterval)
	}
}

func TestLoadConfigInvalidSyncInterval(t *testing.T) {
	setupTestEnv(t)
	os.Setenv("SYNC_INTERVAL_MINUTES", "invalid")
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should fallback to default
	if config.SyncInterval != 60*time.Minute {
		t.Errorf("Expected SyncInterval 60m (default), got %v", config.SyncInterval)
	}
}

func TestLoadConfigZeroSyncInterval(t *testing.T) {
	setupTestEnv(t)
	os.Setenv("SYNC_INTERVAL_MINUTES", "0")
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should fallback to default
	if config.SyncInterval != 60*time.Minute {
		t.Errorf("Expected SyncInterval 60m (default), got %v", config.SyncInterval)
	}
}

func TestLoadConfigNegativeSyncInterval(t *testing.T) {
	setupTestEnv(t)
	os.Setenv("SYNC_INTERVAL_MINUTES", "-10")
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should fallback to default
	if config.SyncInterval != 60*time.Minute {
		t.Errorf("Expected SyncInterval 60m (default), got %v", config.SyncInterval)
	}
}

func TestLoadConfigMissingDatabaseURL(t *testing.T) {
	content := `# Test environment variables without DATABASE_URL
STOCK_API_TOKEN=test-token
STOCK_API_BASE_URL=https://api.test.com
`
	err := os.WriteFile(".env", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	os.Unsetenv("DATABASE_URL")

	config, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if config != nil {
		t.Error("Expected nil config when error occurs")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("Expected error message to contain DATABASE_URL, got %v", err)
	}
}

func TestLoadConfigMissingStockAPIToken(t *testing.T) {
	content := `# Test environment variables without STOCK_API_TOKEN
DATABASE_URL=postgres://test:test@localhost/testdb
STOCK_API_BASE_URL=https://api.test.com
`
	err := os.WriteFile(".env", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	os.Unsetenv("STOCK_API_TOKEN")

	config, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if config != nil {
		t.Error("Expected nil config when error occurs")
	}
	if !strings.Contains(err.Error(), "STOCK_API_TOKEN") {
		t.Errorf("Expected error message to contain STOCK_API_TOKEN, got %v", err)
	}
}

func TestLoadConfigMissingBaseURL(t *testing.T) {
	content := `# Test environment variables without STOCK_API_BASE_URL
DATABASE_URL=postgres://test:test@localhost/testdb
STOCK_API_TOKEN=test-token
`
	err := os.WriteFile(".env", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	os.Unsetenv("STOCK_API_BASE_URL")

	config, err := LoadConfig()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if config != nil {
		t.Error("Expected nil config when error occurs")
	}
	if !strings.Contains(err.Error(), "STOCK_API_BASE_URL") {
		t.Errorf("Expected error message to contain STOCK_API_BASE_URL, got %v", err)
	}
}

func TestLoadConfigDefaultPort(t *testing.T) {
	setupTestEnv(t)
	os.Unsetenv("SERVER_PORT")
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.ServerPort != "8080" {
		t.Errorf("Expected ServerPort 8080 (default), got %s", config.ServerPort)
	}
}

func TestLoadConfigDefaultSyncInterval(t *testing.T) {
	setupTestEnv(t)
	os.Unsetenv("SYNC_INTERVAL_MINUTES")
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.SyncInterval != 60*time.Minute {
		t.Errorf("Expected SyncInterval 60m (default), got %v", config.SyncInterval)
	}
}

func TestGetEnvSyncIntervalValidValues(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected time.Duration
	}{
		{"30 minutes", "30", 30 * time.Minute},
		{"60 minutes", "60", 60 * time.Minute},
		{"120 minutes", "120", 120 * time.Minute},
		{"1 minute", "1", 1 * time.Minute},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SYNC_INTERVAL_MINUTES", tt.value)
			defer os.Unsetenv("SYNC_INTERVAL_MINUTES")

			duration := getEnvSyncInterval()
			if duration != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, duration)
			}
		})
	}
}

func TestGetEnvSyncIntervalInvalidValues(t *testing.T) {
	invalidValues := []string{"", "invalid", "-5", "0", "abc"}

	for _, value := range invalidValues {
		t.Run("invalid_"+value, func(t *testing.T) {
			os.Setenv("SYNC_INTERVAL_MINUTES", value)
			defer os.Unsetenv("SYNC_INTERVAL_MINUTES")

			duration := getEnvSyncInterval()
			expectedDefault := 60 * time.Minute
			if duration != expectedDefault {
				t.Errorf("Expected default %v, got %v", expectedDefault, duration)
			}
		})
	}
}

func TestCORSConfigDefaults(t *testing.T) {
	setupTestEnv(t)
	defer cleanupTestEnv()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	corsConfig := config.CORSSettings

	if len(corsConfig.AllowedOrigins) != 1 || corsConfig.AllowedOrigins[0] != "*" {
		t.Errorf("Expected AllowedOrigins [*], got %v", corsConfig.AllowedOrigins)
	}

	requiredMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	for _, method := range requiredMethods {
		found := false
		for _, m := range corsConfig.AllowedMethods {
			if m == method {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected %s in AllowedMethods", method)
		}
	}

	if len(corsConfig.AllowedHeaders) != 1 || corsConfig.AllowedHeaders[0] != "*" {
		t.Errorf("Expected AllowedHeaders [*], got %v", corsConfig.AllowedHeaders)
	}

	if !corsConfig.AllowCredentials {
		t.Error("Expected AllowCredentials to be true")
	}
}

// Helper functions

func setupTestEnv(t *testing.T) {
	createTestEnvFile(t)

	os.Setenv("DATABASE_URL", "postgres://test:test@localhost/testdb")
	os.Setenv("STOCK_API_TOKEN", "test-token")
	os.Setenv("STOCK_API_BASE_URL", "https://api.test.com")
}

func cleanupTestEnv() {
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("STOCK_API_TOKEN")
	os.Unsetenv("STOCK_API_BASE_URL")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("SYNC_INTERVAL_MINUTES")

	os.Remove(".env")
}

func createTestEnvFile(t *testing.T) {
	content := `# Test environment variables
DATABASE_URL=postgres://test:test@localhost/testdb
STOCK_API_TOKEN=test-token
STOCK_API_BASE_URL=https://api.test.com
`
	err := os.WriteFile(".env", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
}
