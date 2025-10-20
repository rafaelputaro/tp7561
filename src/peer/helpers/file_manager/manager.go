package file_manager

import (
	"fmt"
	"os"
	"tp/common"
	"tp/common/files_common"
	"tp/common/files_common/path_exists"
	"tp/common/files_common/uploader"
	"tp/common/keys"

	"tp/peer/helpers/file_manager/blocks"
	"tp/peer/helpers/file_manager/config_fm"
	"tp/peer/helpers/file_manager/source_file"
	"tp/peer/helpers/file_manager/utils"
)

const MSG_ADD_FILE_FROM_INPUT_DIR = "add file from input directory: %v"
const MSG_ADD_FILE_FROM_UPLOAD_DIR = "add file from upload directory: %v"

type ProcessBlockCallBack func(key []byte, fileName string, data []byte) error

// Leer un archivo por bloques desde el directorio input y enviando los mismos como
// bloques con el formato <blockKey><nextBlockKey><data> los cuales son procesados con
// el callback parámetro
func AddFileFromInputDir(fileName string, processBlock ProcessBlockCallBack) error {
	filePath := utils.GenerateInputFilePath(fileName)
	common.Log.Debugf(MSG_ADD_FILE_FROM_INPUT_DIR, fileName)
	return AddFile(fileName, filePath, processBlock)
}

// Leer un archivo por bloques desde el directorio input y enviando los mismos como
// bloques con el formato <blockKey><nextBlockKey><data> los cuales son procesados con
// el callback parámetro
func AddFileFromUploadDir(fileName string, processBlock ProcessBlockCallBack) error {
	filePath := utils.GenerateIpfsUploadPath(fileName)
	common.Log.Debugf(MSG_ADD_FILE_FROM_UPLOAD_DIR, fileName)
	return AddFile(fileName, filePath, processBlock)
}

// Leer un archivo por bloques desde el path parámetro y enviando los mismos como
// bloques con el formato <blockKey><nextBlockKey><data> los cuales son procesados con
// el callback parámetro
func AddFile(fileName string, filePath string, processBlock ProcessBlockCallBack) error {
	reader, err := source_file.NewFileReader(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()
	// lectura primer bloque
	blkData, blkNum, _, err := reader.Next()
	// si hay error finalizar retornando error
	if err != nil {
		return err
	}
	blkName := blocks.GenerateBlockName(fileName, blkNum)
	blkKey := keys.GetKey(blkName)
	end := false
	for !end {
		// leer un nuevo bloque
		nextBlkData, nextBlkNum, nextEof, nextErr := reader.Next()
		nextName := keys.NULL_KEY_SOURCE_DATA
		// si no hay más bloques el nombre del siguiente es nulo
		if nextEof || nextErr != nil {
			end = true
		} else {
			nextName = blocks.GenerateBlockName(fileName, nextBlkNum)
			blkNum = nextBlkNum
		}
		// key del siguiente bloque
		nextKey := keys.GetKey(nextName)
		// generar bloque
		blockToStore := blocks.GenerateBlockToStore(blkData, blkKey, nextKey)
		// enviar bloque a vecinos
		processBlock(blkKey, blkName, blockToStore)
		// actualizar para el siguiente ciclo
		blkKey = nextKey
		blkName = nextName
		blkData = nextBlkData
	}
	common.Log.Infof(utils.MSG_FILE_ADDED, fileName, blkNum+1)
	return err
}

// Escribe un bloque en un archivo localmente para ser recuperado en una búqueda. Retorna error
// si el archivo ya existe o si se presenta algún error de acceso a disco
func StoreBlock(fileName string, data []byte) error {
	return blocks.StoreBlock(utils.GenerateIpfsStorePath(fileName), data)
}

// Escribe un bloque en un archivo localmente como parte de un archivo a ser recuparado.
// Retorna error si el archivo ya existe o si se presenta algún error de acceso a disco
// Retorna verdadero si es el bloque final, falso caso contrario
func StoreBlockOnDownload(fileName string, data []byte) (bool, error) {
	err := blocks.StoreBlock(utils.GenerateIpfsDownloadPath(fileName), data)
	if err != nil {
		return false, err
	}
	return blocks.IsFinalBlock(data), err
}

// Obtiene un block completo con su header y datos.
func GetBlock(fileName string) ([]byte, error) {
	nBytes, data, err := blocks.ReadBlock(utils.GenerateIpfsStorePath(fileName))
	return data[:nBytes], err
}

// Sube los archivos locales a la red de nodos
func UploadLocalFiles(uploadFile func(fileName string) error) error {
	// leer archivos del directorio
	entries, err := os.ReadDir(config_fm.LocalStorageConfig.InputDataFolder)
	if err != nil {
		common.Log.Errorf(utils.MSG_ERROR_READING_DIRECTORY, err)
	}
	// carga en la red de nodos
	for _, entry := range entries {
		fmt.Printf("- %s", entry.Name())
		if !entry.IsDir() {
			uploadFile(entry.Name())
		}
	}
	return nil
}

// Guarda un archivo en espacio de upload. En caso de ser el último bloque del archivo lo
// reconstruye para luego borrar las partes. <restored><error>
func StoreUploadFilePart(fileName string, part int32, data []byte, endFile bool) (bool, error) {
	var err error = nil
	if part >= 0 {
		err = files_common.StoreFile(utils.GenerateIpfsUploadPartPath(fileName, int(part)), data)
	}
	if err != nil {
		return false, err
	}
	if endFile {
		// Restaurar archivo
		err = uploader.RestoreFile(utils.GenerateIpfsUploadPath(fileName), func(fPart int) string {
			return utils.GenerateIpfsUploadPartPath(fileName, fPart)
		})
		return err == nil, err
	}
	return false, nil
}

// Retorna verdadero si existe el archivo en directorio upload
func FileExistInUpload(fileName string) bool {
	path := utils.GenerateIpfsUploadPath(fileName)
	return path_exists.PathExists(path)
}
