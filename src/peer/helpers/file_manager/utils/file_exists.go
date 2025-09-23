package utils

import (
	"errors"
	"os"
)

// Retorna verdadero si el archivo existe
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return false
	} else {
		return !errors.Is(err, os.ErrNotExist)
	}
}
