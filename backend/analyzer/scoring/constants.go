package scoring

// RatingMapping asigna calificaciones de analistas a valores numéricos
// Ordenadas de mayor (más positivo, 1.0) a menor (más negativo, 0.0)
var RatingMapping = map[string]float64{
	// Recomendaciones muy positivas (0.90-1.00)
	"strong buy": 1.00,
	"strong-buy": 1.00,
	"buy":        0.90,

	// Recomendaciones positivas (0.70-0.89)
	"outperform":        0.80,
	"sector outperform": 0.80,
	"market outperform": 0.80,
	"overweight":        0.75,
	"outperformer":      0.75,
	"positive":          0.75,
	"add":               0.70,

	// Recomendaciones moderadamente positivas (0.60-0.69)
	"accumulate":      0.65,
	"speculative buy": 0.65,
	"moderate buy":    0.60,

	// Recomendaciones neutrales (0.40-0.59)
	"equal weight":   0.50,
	"neutral":        0.50,
	"in-line":        0.50,
	"sector perform": 0.50,
	"sector weight":  0.50,
	"peer perform":   0.50,
	"market perform": 0.45,
	"hold":           0.40,

	// Recomendaciones moderadamente negativas (0.30-0.39)
	"cautious":            0.35,
	"sector underperform": 0.35,
	"underweight":         0.35,
	"underperform":        0.30,

	// Recomendaciones negativas (0.10-0.29)
	"reduce":      0.25,
	"negative":    0.20,
	"weak hold":   0.15,
	"sell":        0.10,
	"strong sell": 0.00,

	// Calificaciones no reconocidas o faltantes
	"unknown": 0.50,
}
