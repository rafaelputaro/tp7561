package files_common

import (
	"fmt"
	"os"
	"tp/common"
	"tp/common/files_common/messages"
)

// Sube los archivos desde un directorio específico
func OpOverDir(pathDir string, opOverFile func(fileName string) error) error {
	// leer archivos del directorio
	entries, err := os.ReadDir(pathDir)
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_READING_DIRECTORY, err)
	}
	// opera sobre los archivos
	for _, entry := range entries {
		fmt.Printf("- %s", entry.Name())
		if !entry.IsDir() {
			opOverFile(entry.Name())
		}
	}
	return nil
}
