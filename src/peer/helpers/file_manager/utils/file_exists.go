package utils

import (
	"os"
)

// Retorna verdadero si la ruta existe
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
