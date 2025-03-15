// Package database proporciona funcionalidades para interactuar con la base de datos.
package database

import (
	"github.com/JhonierSerna14/STOCK-VIZ/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB inicializa la conexión a la base de datos y realiza las migraciones necesarias
func InitDB(dsn string) (*gorm.DB, error) {
	// Inicializar conexión a la base de datos con GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Configurar el nivel de logging para mostrar solo errores
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	// Migración automática del esquema de base de datos
	if err := db.AutoMigrate(&models.Stock{}); err != nil {
		return nil, err
	}

	return db, nil
}
