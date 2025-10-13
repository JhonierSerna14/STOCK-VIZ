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
}

// ActionSentimentMapping asigna acciones de brokers a valores de sentimiento
// Valores: 1.0 (positivo), 0.0 (neutral), -1.0 (negativo)
var ActionSentimentMapping = map[string]float64{
	// Acciones positivas
	"upgraded":         1.0,
	"upgraded by":      1.0,
	"target raised":    1.0,
	"target raised by": 1.0,
	"initiated":        1.0,
	"initiated by":     1.0,
	"target set":       0.5,
	"target set by":    0.5,
	"raised":           1.0,
	"improve":          1.0,

	// Acciones neutrales
	"reiterated":    0.0,
	"reiterated by": 0.0,
	"maintained":    0.0,
	"confirmed":     0.0,
	"unchanged":     0.0,

	// Acciones negativas
	"downgraded":        -1.0,
	"downgraded by":     -1.0,
	"target lowered":    -1.0,
	"target lowered by": -1.0,
	"lowered":           -1.0,
	"reduced":           -1.0,
	"suspended":         -1.0,
}
