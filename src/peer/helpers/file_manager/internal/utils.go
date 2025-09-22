package internal

import (
	"errors"
	"os"
	"tp/peer/helpers"
)

const MSG_FILE_COULD_NOT_BE_OPONED = "the file could not be opened: %v"
const MSG_ERROR_CLOSING_FILE = "error closing file: %v"
const MSG_ERROR_READING_FILE = "error reading file: %v"
const INVALID_BLOCK_NUMBER = -1
const BLOCK_SIZE = 256 * 1024 // tamaño de los bloques en bytes
const HEADER_BLOCK_FILE_SIZE = 2 * helpers.LENGTH_KEY_IN_BITS
const MAX_BLOCK_FILE_SIZE = HEADER_BLOCK_FILE_SIZE + BLOCK_SIZE

// Retorna verdadero si el archivo existe
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err == nil {
		return false
	} else {
		return !errors.Is(err, os.ErrNotExist)
	}
}
