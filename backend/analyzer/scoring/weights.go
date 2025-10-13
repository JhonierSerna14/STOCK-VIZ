package scoring

// FactorWeights define los pesos para cada factor de análisis
type FactorWeights struct {
	Rating      float64 // Peso para la calificación de analistas con tendencias históricas
	Target      float64 // Peso para el cambio de precio objetivo
	TargetTrend float64 // Peso para la evolución histórica de precios objetivo
	Volatility  float64 // Peso para la consistencia histórica (menor volatilidad = mayor confianza)
	Recency     float64 // Peso para la actualidad de la información
}

// DefaultWeights proporciona los pesos recalibrados con análisis temporal
var DefaultWeights = FactorWeights{
	Rating:      0.30, // Reducido para dar espacio a tendencias
	Target:      0.25, // Reducido para dar espacio a tendencias de target
	TargetTrend: 0.20, // Nuevo: tendencia histórica de precios objetivo
	Volatility:  0.15, // Nuevo: consistencia de recomendaciones
	Recency:     0.10, // Mantiene importancia de actualidad
}
