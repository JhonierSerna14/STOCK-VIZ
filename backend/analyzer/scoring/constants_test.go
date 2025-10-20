package scoring

import (
	"testing"
)

func TestRatingMappingPositiveRatings(t *testing.T) {
	// Test positive ratings have high scores
	positiveRatings := []string{"strong buy", "strong-buy", "buy", "outperform", "overweight"}

	for _, rating := range positiveRatings {
		score, exists := RatingMapping[rating]
		if !exists {
			t.Errorf("Rating %s should exist in mapping", rating)
			continue
		}
		if score < 0.70 {
			t.Errorf("Rating %s should have score >= 0.70, got %f", rating, score)
		}
	}
}

func TestRatingMappingNegativeRatings(t *testing.T) {
	// Test negative ratings have low scores
	negativeRatings := []string{"sell", "strong sell", "underperform", "reduce"}

	for _, rating := range negativeRatings {
		score, exists := RatingMapping[rating]
		if !exists {
			t.Errorf("Rating %s should exist in mapping", rating)
			continue
		}
		if score > 0.30 {
			t.Errorf("Rating %s should have score <= 0.30, got %f", rating, score)
		}
	}
}

func TestRatingMappingNeutralRatings(t *testing.T) {
	// Test neutral ratings have medium scores
	neutralRatings := []string{"neutral", "hold", "equal weight", "market perform"}

	for _, rating := range neutralRatings {
		score, exists := RatingMapping[rating]
		if !exists {
			t.Errorf("Rating %s should exist in mapping", rating)
			continue
		}
		if score < 0.40 {
			t.Errorf("Rating %s should have score >= 0.40, got %f", rating, score)
		}
		if score > 0.59 {
			t.Errorf("Rating %s should have score <= 0.59, got %f", rating, score)
		}
	}
}

func TestRatingMappingSpecificValues(t *testing.T) {
	// Test specific rating values
	tests := []struct {
		rating   string
		expected float64
	}{
		{"strong buy", 1.00},
		{"buy", 0.90},
		{"neutral", 0.50},
		{"hold", 0.40},
		{"sell", 0.10},
		{"strong sell", 0.00},
		{"unknown", 0.50},
	}

	for _, tt := range tests {
		t.Run(tt.rating, func(t *testing.T) {
			score, exists := RatingMapping[tt.rating]
			if !exists {
				t.Errorf("Rating %s should exist", tt.rating)
				return
			}
			if score != tt.expected {
				t.Errorf("Rating %s should have score %f, got %f", tt.rating, tt.expected, score)
			}
		})
	}
}

func TestRatingMappingAllValuesValid(t *testing.T) {
	// Test that all values are between 0 and 1
	for rating, score := range RatingMapping {
		if score < 0.0 {
			t.Errorf("Rating %s score should be >= 0, got %f", rating, score)
		}
		if score > 1.0 {
			t.Errorf("Rating %s score should be <= 1, got %f", rating, score)
		}
	}
}

func TestRatingMappingCaseVariations(t *testing.T) {
	// Test case variations that should exist
	caseTests := []string{
		"strong buy",
		"strong-buy",
		"outperform",
		"overweight",
		"equal weight",
		"market perform",
		"sector perform",
		"strong sell",
	}

	for _, rating := range caseTests {
		_, exists := RatingMapping[rating]
		if !exists {
			t.Errorf("Rating %s should exist in mapping", rating)
		}
	}
}

func TestRatingMappingUnknownRating(t *testing.T) {
	// Test unknown rating fallback
	score, exists := RatingMapping["unknown"]
	if !exists {
		t.Error("Unknown rating should exist")
	}
	if score != 0.50 {
		t.Errorf("Unknown rating should have score 0.50 (neutral), got %f", score)
	}
}

func TestRatingMappingCountEntries(t *testing.T) {
	// Verify we have reasonable number of entries
	if len(RatingMapping) <= 20 {
		t.Errorf("Should have sufficient rating mappings, got %d entries", len(RatingMapping))
	}
}
