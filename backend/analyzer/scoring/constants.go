package scoring

// RatingMapping asigna calificaciones de analistas a valores numéricos
// Ordenadas de mayor (más positivo, 1.0) a menor (más negativo, 0.0)
var RatingMapping = map[string]float64{
	// Recomendaciones muy positivas (0.90-1.00)
	"strong-buy": 1.00,

	// Recomendaciones positivas (0.70-0.89)
	"buy":               0.80,
	"outperform":        0.75,
	"sector outperform": 0.75,
	"market outperform": 0.75,
	"overweight":        0.70,
	"outperformer":      0.70,
	"positive":          0.70,

	// Recomendaciones moderadamente positivas (0.60-0.69)
	"add":             0.65,
	"accumulate":      0.65,
	"speculative buy": 0.65,

	// Recomendaciones neutrales (0.40-0.59)
	"equal weight":   0.50,
	"neutral":        0.50,
	"in-line":        0.50,
	"sector perform": 0.50,
	"sector weight":  0.50,
	"peer perform":   0.50,
	"market perform": 0.45,
	"hold":           0.40,

	// Recomendaciones negativas (0.20-0.39)
	"sector underperform": 0.30,
	"underweight":         0.30,
	"underperform":        0.30,
	"reduce":              0.30,
	"negative":            0.30,

	// Recomendaciones muy negativas (0.00-0.19)
	"sell": 0.00,
}
