package files_common

import (
	"os"
	"tp/common"
	"tp/common/files_common/messages"
	"tp/common/files_common/path_exists"
)

// Escribe un bloque de datos en un archivo localmente. Retorna error si el archivo ya existe o si se
// presenta algún error de acceso a disco
func StoreFile(filepath string, data []byte) error {
	// chequear si el archivo ya existe
	if path_exists.PathExists(filepath) {
		common.Log.Debugf(messages.MSG_ERROR_FILE_EXIST)
		return os.ErrExist
	}
	file, err := os.Create(filepath)
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_CREATING_FILE, err)
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		common.Log.Errorf(messages.MSG_ERROR_WRITING_FILE, err)
		return err
	}
	return nil
}
