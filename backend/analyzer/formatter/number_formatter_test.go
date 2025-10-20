package formatter

import (
	"testing"
)

func TestNumberFormatterExtractNumber(t *testing.T) {
	formatter := &NumberFormatter{}

	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{
			name:     "Simple number",
			input:    "123.45",
			expected: 123.45,
			hasError: false,
		},
		{
			name:     "Number with dollar sign",
			input:    "$123.45",
			expected: 123.45,
			hasError: false,
		},
		{
			name:     "Number with commas",
			input:    "1,234.56",
			expected: 1234.56,
			hasError: false,
		},
		{
			name:     "Number with dollar and commas",
			input:    "$1,234.56",
			expected: 1234.56,
			hasError: false,
		},
		{
			name:     "Number with spaces",
			input:    " 123.45 ",
			expected: 123.45,
			hasError: false,
		},
		{
			name:     "Complex formatting",
			input:    " $1,234,567.89 ",
			expected: 1234567.89,
			hasError: false,
		},
		{
			name:     "Zero value",
			input:    "0",
			expected: 0,
			hasError: false,
		},
		{
			name:     "Zero with dollar",
			input:    "$0.00",
			expected: 0,
			hasError: false,
		},
		{
			name:     "Integer value",
			input:    "150",
			expected: 150,
			hasError: false,
		},
		{
			name:     "Integer with dollar",
			input:    "$150",
			expected: 150,
			hasError: false,
		},
		{
			name:     "Large number",
			input:    "$10,000,000.00",
			expected: 10000000,
			hasError: false,
		},
		{
			name:     "Invalid string",
			input:    "invalid",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Only dollar sign",
			input:    "$",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Mixed invalid characters",
			input:    "$abc123.45",
			expected: 0,
			hasError: true,
		},
		{
			name:     "Multiple decimal points",
			input:    "123.45.67",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatter.ExtractNumber(tt.input)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
				if result != 0.0 {
					t.Errorf("Expected result 0.0 on error, got %f", result)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %f, got %f", tt.expected, result)
				}
			}
		})
	}
}

func TestNumberFormatterExtractNumberEdgeCases(t *testing.T) {
	formatter := &NumberFormatter{}

	// Test with very large numbers
	result, err := formatter.ExtractNumber("999999999999.99")
	if err != nil {
		t.Errorf("Expected no error for large number, got %v", err)
	}
	if result != 999999999999.99 {
		t.Errorf("Expected 999999999999.99, got %f", result)
	}

	// Test with very small numbers
	result, err = formatter.ExtractNumber("0.001")
	if err != nil {
		t.Errorf("Expected no error for small number, got %v", err)
	}
	if result != 0.001 {
		t.Errorf("Expected 0.001, got %f", result)
	}

	// Test with negative numbers
	result, err = formatter.ExtractNumber("-123.45")
	if err != nil {
		t.Errorf("Expected no error for negative number, got %v", err)
	}
	if result != -123.45 {
		t.Errorf("Expected -123.45, got %f", result)
	}

	// Test with negative dollar amounts
	result, err = formatter.ExtractNumber("-$123.45")
	if err != nil {
		t.Errorf("Expected no error for negative dollar amount, got %v", err)
	}
	if result != -123.45 {
		t.Errorf("Expected -123.45, got %f", result)
	}
}
