// Package models contiene las estructuras de datos y modelos utilizados en la aplicación.
package models

import "time"

// Stock representa información sobre una recomendación o acción de una acción bursátil.
// Se utiliza tanto para la serialización JSON como para el mapeo con la base de datos mediante GORM.
type Stock struct {
	// Ticker es el símbolo único que identifica la acción en el mercado bursátil
	// Forma parte de la clave primaria compuesta junto con Time
	Ticker string `json:"ticker" gorm:"primaryKey"`

	// Company es el nombre de la empresa a la que pertenece la acción
	Company string `json:"company"`

	// Brokerage es el nombre de la casa de bolsa o entidad que emite la recomendación
	Brokerage string `json:"brokerage"`

	// Action representa el tipo de acción o recomendación (ej: "comprar", "vender", "mantener")
	Action string `json:"action"`

	// RatingFrom indica la calificación anterior de la acción
	RatingFrom string `json:"rating_from"`

	// RatingTo indica la nueva calificación asignada a la acción
	RatingTo string `json:"rating_to"`

	// TargetFrom indica el precio objetivo anterior
	TargetFrom string `json:"target_from"`

	// TargetTo indica el nuevo precio objetivo
	TargetTo string `json:"target_to"`

	// Time es la fecha y hora en que se emitió la recomendación
	// Forma parte de la clave primaria compuesta junto con Ticker
	Time time.Time `json:"time" gorm:"primaryKey"`
}
