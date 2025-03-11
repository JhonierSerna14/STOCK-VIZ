package utils

import "strings"

// NormalizeAction estandariza el texto de una acción
func NormalizeAction(action string) string {
	norm := strings.ToLower(strings.TrimSpace(action))
	norm = strings.TrimSuffix(norm, " by")
	return norm
}
