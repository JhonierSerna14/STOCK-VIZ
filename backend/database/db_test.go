package database

import (
	"testing"
)

func TestInitDBInvalidDSN(t *testing.T) {
	// Test with invalid DSN
	db, err := InitDB("invalid://invalid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if db != nil {
		t.Error("Expected nil db")
	}
}

func TestInitDBPostgresDSN(t *testing.T) {
	// Test with invalid PostgreSQL DSN (should fail gracefully)
	db, err := InitDB("postgres://invalid:invalid@localhost/invalid")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if db != nil {
		t.Error("Expected nil db")
	}
}

func TestInitDBEmptyDSN(t *testing.T) {
	// Test with empty DSN
	db, err := InitDB("")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if db != nil {
		t.Error("Expected nil db")
	}
}

func TestInitDBMalformedDSN(t *testing.T) {
	// Test with malformed DSN
	db, err := InitDB("not-a-valid-dsn")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if db != nil {
		t.Error("Expected nil db")
	}
}
