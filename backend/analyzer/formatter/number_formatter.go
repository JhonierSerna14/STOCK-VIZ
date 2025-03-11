package formatter

import (
	"fmt"
	"strconv"
	"strings"
)

// NumberFormatter proporciona funciones para el formateo de números
type NumberFormatter struct{}

// ExtractNumber convierte una cadena que representa un valor monetario a float64
func (f *NumberFormatter) ExtractNumber(value string) (float64, error) {
	clean := strings.ReplaceAll(value, ",", "")
	clean = strings.TrimSpace(strings.Replace(clean, "$", "", -1))
	num, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting value %s: %w", value, err)
	}
	return num, nil
}

// NormalizeAction estandariza el formato de las acciones
func (f *NumberFormatter) NormalizeAction(action string) string {
	norm := strings.ToLower(strings.TrimSpace(action))
	return strings.TrimSuffix(norm, " by")
}
