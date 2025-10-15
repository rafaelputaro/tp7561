package uploader

import (
	"io"
	"os"
	"tp/common"
	"tp/common/files_common/messages"
)

// Reconstruye un archivo en base a sus partes y lo almacena en el archivo de salida
func RestoreFile(outputFilePath string, generateNextPartPath func(int) string) error {
	// creo archivo de salida
	file, err := os.Create(outputFilePath)
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_CREATING_FILE, err)
		return err
	}
	defer file.Close()
	// iniciar recuperación
	part := 0
	for {
		path := generateNextPartPath(part)
		// leer siguiente bloque
		numBytes, data, err := ReadPart(path)
		if err != nil {
			common.Log.Errorf(messages.MSG_ERROR_READING_PART, err)
			return err
		}
		if numBytes <= 0 {
			break
		}
		// escribir en archivo de recuperación
		if _, err := file.Write(data); err != nil {
			common.Log.Errorf(messages.MSG_ERROR_WRITING_FILE, err)
			return err
		}
		part++
	}
	common.Log.Infof(messages.MSG_FILE_RESTORED, outputFilePath)
	return nil
}

// Lee un archivo completo
func ReadPart(path string) (int, []byte, error) {
	// abrir archivo
	f, err := os.Open(path)
	if err != nil {
		common.Log.Errorf(messages.MSG_FILE_COULD_NOT_BE_OPONED, err)
		return 0, nil, err
	}
	defer f.Close()
	fileContent := make([]byte, MAX_PART_SIZE)
	// leer bloque completo
	nBytes, err := f.Read(fileContent)
	if err != nil {
		if err != io.EOF {
			common.Log.Errorf(messages.MSG_ERROR_READING_FILE, err)
			return nBytes, nil, err
		}
	}
	return nBytes, fileContent[0:nBytes], nil
}
