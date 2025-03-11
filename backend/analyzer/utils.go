// Contiene funciones auxiliares
package analyzer

import (
	"fmt"
	"strconv"
	"strings"
)

// extractNumber convierte una cadena que representa un valor monetario (ej: "$55.00")
// a su equivalente numérico (float64). Maneja la limpieza de símbolos monetarios
// y separadores de miles.
func (a *StockAnalyzer) extractNumber(value string) (float64, error) {
	// Elimina las comas y el símbolo "$"
	clean := strings.ReplaceAll(value, ",", "")
	clean = strings.TrimSpace(strings.Replace(clean, "$", "", -1))
	num, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0, fmt.Errorf("error converting value %s: %w", value, err)
	}
	return num, nil
}

// normalizeAction estandariza el texto de una acción eliminando variaciones
// no significativas y convirtiéndolo a un formato consistente para su análisis.
func normalizeAction(action string) string {
	norm := strings.ToLower(strings.TrimSpace(action))
	norm = strings.TrimSuffix(norm, " by")
	return norm
}
