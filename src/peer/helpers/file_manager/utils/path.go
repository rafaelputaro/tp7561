package utils

import (
	"os"
	"strconv"
	"tp/common"
	"tp/common/files_common/path_exists"
	"tp/peer/helpers/file_manager/config_fm"
)

// Retorna el path completo de un archivo de entrada de acuerdo a la configuración
func GenerateInputFilePath(fileName string) string {
	return config_fm.LocalStorageConfig.InputDataFolder + "/" + fileName
}

// Retorna el path completo del archivo de acuerdo a la configuración
func GenerateIpfsStorePath(fileName string) string {
	return config_fm.LocalStorageConfig.StoreIpfsFolder + "/" + fileName
}

// Retorna el path completo de un archivo recuperado <directory restore>/<fileName>
func GenerateIpfsRestorePath(fileName string) string {
	return config_fm.LocalStorageConfig.RestoreIpfsFolder + "/" + fileName
}

// Retorna el path completo de un archivo situado en la carpeta de descargas
// <directory down>/<filename>
func GenerateIpfsDownloadPath(fileName string) string {
	return config_fm.LocalStorageConfig.DownloadIpfsFolder + "/" + fileName
}

// Retorna el path completo de un archivo situado en la carpeta de upload
// <directory down>/<filename>
func GenerateIpfsUploadPath(fileName string) string {
	return config_fm.LocalStorageConfig.UploadIpfsFolder + "/" + fileName
}

// Retorna el path completo de una parte de un archivo descargado de la red de nodos
// <directory down>/<filename>.part<blockNumber>
func GenerateIpfsDownloadPartPath(fileName string, blockNumber int) string {
	return GenerateIpfsDownloadPath(fileName + GeneratePartExtension(blockNumber))
}

// Retorna el path completo de una parte de un archivo subido desde fuera de la red de nodos
// <directory upload>/<filename>.part<blockNumber>
func GenerateIpfsUploadPartPath(fileName string, blockNumber int) string {
	return GenerateIpfsUploadPath(fileName + GeneratePartExtension(blockNumber))
}

// Retorna la extensión de la parte de un archivo
func GeneratePartExtension(blockNumber int) string {
	if blockNumber > 0 {
		return ".part" + strconv.Itoa(blockNumber)
	}
	return ""
}

// Lee las variables de entorno que establecen la configuración de almacenamiento and clean the store
func InitStore() {
	config_fm.LoadConfig()
	// limpiar archivos viejos
	CleanStore()
	// crear carpetas necesarias
	CreateStoreFolders()
}

// Limpia el store
func CleanStore() {
	path := GenerateIpfsStorePath("")
	if path_exists.PathExists(path) {
		err := os.RemoveAll(path)
		if err != nil {
			common.Log.Errorf(MSG_ERROR_ON_CLEAN_STORE, path, err)
			return
		}
		common.Log.Infof(MSG_STORE_HAS_BENN_CLEANED)
	}
}

func CreateStoreFolders() {
	// crear store folder
	path := GenerateIpfsStorePath("")
	err := os.Mkdir(path, 0755)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_CREATING_FOLDER, err)
	}
	// crear restore folder dentro de store
	path = GenerateIpfsRestorePath("")
	err = os.Mkdir(path, 0755)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_CREATING_FOLDER, err)
	}
	// crear download folder dentro de store
	path = GenerateIpfsDownloadPath("")
	err = os.Mkdir(path, 0755)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_CREATING_FOLDER, err)
	}
	// crear upload folder dentro de store
	path = GenerateIpfsUploadPath("")
	err = os.Mkdir(path, 0755)
	if err != nil {
		common.Log.Errorf(MSG_ERROR_CREATING_FOLDER, err)
	}
}
