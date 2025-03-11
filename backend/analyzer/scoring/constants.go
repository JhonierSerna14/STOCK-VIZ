package scoring

// RatingMapping asigna calificaciones de analistas a valores numéricos
var RatingMapping = map[string]float64{
	"strong buy":          1.0,
	"buy":                 0.8,
	"outperform":          0.75,
	"overweight":          0.70,
	"equal weight":        0.5,
	"neutral":             0.5,
	"in-line":             0.5,
	"hold":                0.4,
	"market perform":      0.45,
	"sector outperform":   0.75,
	"sector perform":      0.5,
	"sector underperform": 0.3,
	"underweight":         0.3,
	"underperform":        0.3,
	"sell":                0.0,
}
