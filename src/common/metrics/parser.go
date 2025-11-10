package common_metrics

import (
	"regexp"
	"strconv"
)

// Retorna el segundo número que aparece en el nombre de un archivo
func ParseFileNumber(fileName string) float64 {
	patron := `(\d+)`
	ren := regexp.MustCompile(patron)
	found := ren.FindAllString(fileName, -1)
	if len(found) > 1 {
		converted, err := strconv.Atoi(found[1])
		if err == nil {
			return float64(converted)
		}
	}
	return -1
}
