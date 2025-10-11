package blocks

import (
	"errors"
	"fmt"
	"io"
	"os"
	"tp/common"
	"tp/common/files_common"
	"tp/common/files_common/messages"
	"tp/peer/helpers"
	"tp/peer/helpers/file_manager/config_fm"
	"tp/peer/helpers/file_manager/utils"
)

// Lee un bloque del espacio local y lo retorna <nro bytes leídos><data><error>.
// En caso de no encontrar el archivo o que se encuentre mal formateado retorna error
func ReadBlock(path string) (int, []byte, error) {
	// abrir archivo
	f, err := os.Open(path)
	if err != nil {
		common.Log.Errorf(messages.MSG_FILE_COULD_NOT_BE_OPONED, err)
		return 0, nil, err
	}
	defer f.Close()
	fileContent := make([]byte, config_fm.MAX_BLOCK_FILE_SIZE)
	// leer bloque completo
	nBytes, err := f.Read(fileContent)
	if err != nil {
		if err != io.EOF {
			common.Log.Errorf(messages.MSG_ERROR_READING_FILE, err)
			return nBytes, nil, err
		}
	}
	return nBytes, fileContent, nil
}

// Lee un bloque del espacio local y lo retorna <blockKey>, <nextBlockKey>, <data> <error>.
// En caso de no encontrar el archivo o que se encuentre mal formateado retorna error
func ReadAndParseBlock(path string) ([]byte, []byte, []byte, error) {
	// leer archivo completo
	nBytes, fileContent, err := ReadBlock(path)
	if err != nil {
		return nil, nil, nil, err
	}
	// parseo de contenido del archivo
	if nBytes < config_fm.HEADER_BLOCK_FILE_SIZE {
		msg := fmt.Sprintf(utils.MSG_ERROR_READING_HEADER, path)
		return nil, nil, nil, errors.New(msg)
	}
	// parseo header
	blockKey := fileContent[:helpers.LENGTH_KEY_IN_BITS]
	nextBlockKey := fileContent[helpers.LENGTH_KEY_IN_BYTES:config_fm.HEADER_BLOCK_FILE_SIZE]
	// parseo data
	data := []byte{}
	if nBytes-config_fm.HEADER_BLOCK_FILE_SIZE > 0 {
		data = fileContent[config_fm.HEADER_BLOCK_FILE_SIZE:nBytes]
	}
	return blockKey, nextBlockKey, data, nil
}

// Lee un bloque del espacio local y lo retorna <blockKey>, <nextBlockKey>, <data>. Luego de leer
// el bloque borra el archivo. En caso de no encontrar el archivo o que se encuentre mal formateado
// retorna error
func ReadAndDeleteBlock(path string) ([]byte, []byte, []byte, error) {
	blockKey, nextBlockKey, data, err := ReadAndParseBlock(path)
	if err != nil {
		// borrar archivo
		if errRem := os.Remove(path); errRem != nil {
			common.Log.Errorf(utils.MSG_ERROR_ON_DELETE_FILE, errRem)
			return nil, nil, nil, errRem
		}
	}
	return blockKey, nextBlockKey, data, err
}

// Reconstruye un archivo en base a sus partes y lo almacena en la carpeta recovered retornando el
// el path de dicho archivo
func RestoreFile(fileName string) (string, error) {
	// creo archivo de salida
	outputFile := utils.GenerateIpfsRestorePath(fileName)
	file, err := os.Create(outputFile)
	if err != nil {
		common.Log.Errorf(utils.MSG_ERROR_CREATING_FILE, err)
		return "", err
	}
	defer file.Close()
	// iniciar recuperación
	blockNumber := 0
	for {
		path := utils.GenerateIpfsDownloadPartPath(fileName, blockNumber)
		// leer siguiente bloque
		_, nextBlockKey, data, err := ReadAndParseBlock(path)
		if err != nil {
			common.Log.Errorf(utils.MSG_ERROR_READING_BLOCK, err)
			return "", err
		}
		// escribir en archivo de recuperación
		if _, err := file.Write(data); err != nil {
			common.Log.Errorf(messages.MSG_ERROR_WRITING_FILE, err)
			return "", err
		}
		// si la clave siguiente es nula es el último bloque
		if helpers.IsNullKey(nextBlockKey) {
			break
		}
		blockNumber++
	}
	common.Log.Infof(messages.MSG_FILE_RESTORED, outputFile)
	return outputFile, nil
}

// Retorna el nombre del bloque en base al nombre del archivo y el número de bloque
func GenerateBlockName(fileName string, blockNumber int) string {
	toReturn := fileName
	if blockNumber != 0 {
		toReturn += utils.GeneratePartExtension(blockNumber)
	}
	return toReturn
}

// Construye un bloque de la siguiente manera <blockKey><nextBlockKey><data>
func GenerateBlockToStore(data []byte, key []byte, nextKey []byte) []byte {
	block := []byte{}
	block = append(block, key...)
	block = append(block, nextKey...)
	block = append(block, data...)
	return block
}

// Escribe un bloque en un archivo localmente. Retorna error si el archivo ya existe o si se
// presenta algún error de acceso a disco
func StoreBlock(filepath string, data []byte) error {
	return files_common.StoreFile(filepath, data)
	/*
		// chequear si el archivo ya existe
		if path_exists.PathExists(filepath) {
			common.Log.Debugf(utils.MSG_ERROR_FILE_EXIST)
			return os.ErrExist
		}
		file, err := os.Create(filepath)
		if err != nil {
			common.Log.Errorf(utils.MSG_ERROR_CREATING_FILE, err)
			return err
		}
		defer file.Close()
		_, err = file.Write(data)
		if err != nil {
			common.Log.Errorf(utils.MSG_ERROR_WRITING_FILE, err)
			return err
		}
		return nil
	*/
}

// Retorna verdadero si el dato contiene el bloque final
func IsFinalBlock(data []byte) bool {
	return helpers.IsNullKey(GetNextBlock(data))
}

// Obtiene el id del siguiente bloque desde un dato con header
func GetNextBlock(data []byte) []byte {
	length := len(data)
	if length < config_fm.HEADER_BLOCK_FILE_SIZE {
		common.Log.Debugf(utils.MSG_ERROR_HEADER_SIZE, length)
		return helpers.GetNullKey()
	}
	return data[helpers.LENGTH_KEY_IN_BYTES:config_fm.HEADER_BLOCK_FILE_SIZE]
}
