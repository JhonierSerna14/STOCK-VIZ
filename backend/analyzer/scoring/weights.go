package scoring

// FactorWeights define los pesos para cada factor de análisis
type FactorWeights struct {
	Rating  float64
	Target  float64
	Broker  float64
	Recency float64
}

// DefaultWeights proporciona los pesos predeterminados
var DefaultWeights = FactorWeights{
	Rating:  0.25,
	Target:  0.25,
	Broker:  0.30,
	Recency: 0.20,
}
