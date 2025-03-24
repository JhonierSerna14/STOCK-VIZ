package scoring

// FactorWeights define los pesos para cada factor de análisis
type FactorWeights struct {
	Rating  float64 // Peso para la calificación de analistas
	Target  float64 // Peso para el cambio de precio objetivo
	Broker  float64 // Peso para el consenso de brokers
	Recency float64 // Peso para la actualidad de la información
}

// DefaultWeights proporciona los pesos predeterminados
var DefaultWeights = FactorWeights{
	// Asignamos mayor peso a los factores más importantes y ordenados por relevancia
	Rating:  0.35, // La calificación de analistas es el factor más determinante
	Target:  0.30, // El precio objetivo es el segundo factor más importante
	Broker:  0.25, // El consenso de brokers tiene una importancia media
	Recency: 0.10, // La actualidad tiene menor peso pero sigue siendo relevante
}
