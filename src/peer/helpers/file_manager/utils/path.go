package utils

import (
	"strconv"
	"tp/peer/helpers/file_manager/config_fm"
)

// Retorna el path completo del archivo de acuerdo a la configuración
func GenerateInputFilePath(fileName string) string {
	return config_fm.LocalStorageConfig.InputDataFolder + "/" + fileName
}

// Retorna el path completo del archivo de acuerdo a la configuración
func GenerateIpfsFilePath(fileName string) string {
	return config_fm.LocalStorageConfig.StoreIpfsFolder + "/" + fileName
}

// Retorna el path de la carpeta recover
func GenerateRecoverFolderPath() string {
	return config_fm.LocalStorageConfig.StoreIpfsFolder + "/" + config_fm.RECOVERED_FOLDER
}

// Retorna el path completo de un archivo recuperado <folder rec>/<fileName>.rec
func GenerateRecoverPath(fileName string) string {
	return GenerateRecoverFolderPath() + "/" + fileName + config_fm.RECOVERED_EXTENSION
}

// Retorna el path completo de una parte de un archivo descargado de la red de nodos
func GenertaIpfsDownloadPartPath(fileName string, blockNumber int) string {
	return GenerateRecoverFolderPath() + "/" + fileName + GeneratePartExtension(blockNumber)
}

// Retorna el path completo de una parte de un archivo descargado de la red de nodos
func GenerateIpfsDownloadPath(fileName string) string {
	return GenerateRecoverPath(fileName + generateDownloadPathExtension())
}

// Agrega al archivo la extensión de descarga de una parte
func generateDownloadPathExtension() string {
	return ".down"
}

func GeneratePartExtension(blockNumber int) string {
	if blockNumber > 0 {
		return ".part" + strconv.Itoa(blockNumber)
	}
	return ""
}
