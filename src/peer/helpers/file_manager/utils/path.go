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
func GenertaIpfsRecoverFolderPath() string {
	return config_fm.LocalStorageConfig.StoreIpfsFolder + "/" + config_fm.RECOVERED_FOLDER
}

// Retorna el path completo de un archivo recuperado
func GenertaIpfsRecoverPath(fileName string) string {
	return GenertaIpfsRecoverFolderPath() + "/" + fileName + ".rec"
}

// Retorna el path completo de una parte de un archivo descargado de la red de nodos
func GenertaIpfsDownloadPartPath(fileName string, blockNumber int) string {
	return GenertaIpfsRecoverPath(fileName + GeneratePartExtension(blockNumber))
}

func GeneratePartExtension(blockNumber int) string {
	if blockNumber > 0 {
		return ".part" + strconv.Itoa(blockNumber)
	}
	return ""
}
