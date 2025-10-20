package helpers

import (
	"os"
	"strconv"
	"tp/common"
	"tp/common/files_common/path_exists"
)

const MSG_ERROR_CREATING_FOLDER = "error creating folder: %v"
const MSG_ERROR_ON_CLEAN_STORE = "error cleaning the store: %v | error: %v"
const MSG_STORE_HAS_BEEN_CLEANED = "the store has been cleaned"

// Retorna el path completo de un archivo de entrada de acuerdo a la configuración
func GenerateInputFilePath(config Config, fileName string) string {
	return config.InputDataFolder + "/" + fileName
}

// Retorna el path completo del archivo de acuerdo a la configuración
func GenerateStorePath(config Config, fileName string) string {
	return config.StoreFolder + "/" + fileName
}

// Retorna el path completo de un archivo situado en la carpeta de descargas
// <directory down>/<filename>
func GenerateDownloadPath(config Config, fileName string) string {
	return config.DownloadFolder + "/" + fileName
}

// Retorna el path completo de una parte de un archivo descargado de la red de nodos
// <directory down>/<filename>.part<blockNumber>
func GenerateDownloadPartPath(config Config, fileName string, blockNumber int) string {
	return GenerateDownloadPath(config, fileName+GeneratePartExtension(blockNumber))
}

// Retorna la extensión de la parte de un archivo
func GeneratePartExtension(blockNumber int) string {
	if blockNumber > 0 {
		return ".part" + strconv.Itoa(blockNumber)
	}
	return ""
}

// Lee las variables de entorno que establecen la configuración de almacenamiento and clean the store
func InitStore(config Config) {
	// limpiar archivos viejos
	CleanStore(config)
	// crear carpetas necesarias
	CreateStoreFolders(config)
}

// Limpia el store
func CleanStore(config Config) {
	path := GenerateStorePath(config, "")
	if path_exists.PathExists(path) {
		err := os.RemoveAll(path)
		if err != nil {
			common.Log.Errorf(MSG_ERROR_ON_CLEAN_STORE, path, err)
			return
		}
		common.Log.Infof(MSG_STORE_HAS_BEEN_CLEANED)
	}
}

func CreateStoreFolders(config Config) {
	// crear store folder
	path := GenerateStorePath(config, "")
	err := os.Mkdir(path, 0755)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_CREATING_FOLDER, err)
	}
	// crear download folder dentro de store
	path = GenerateDownloadPath(config, "")
	err = os.Mkdir(path, 0755)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_CREATING_FOLDER, err)
	}
}
