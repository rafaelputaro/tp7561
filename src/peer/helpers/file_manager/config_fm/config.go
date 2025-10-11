package config_fm

import (
	"os"
	"tp/common"
)

type StorageConfig struct {
	InputDataFolder    string
	StoreIpfsFolder    string
	DownloadIpfsFolder string
	RestoreIpfsFolder  string
	UploadIpfsFolder   string
}

var LocalStorageConfig StorageConfig

// Retorna una nueva instancia de la configuración
func NewStorageConfig(inputDataFolder string, storeIpfsFolder string) *StorageConfig {
	config := &StorageConfig{
		InputDataFolder:    inputDataFolder,
		StoreIpfsFolder:    storeIpfsFolder,
		DownloadIpfsFolder: generateDownloadIpfsFolder(storeIpfsFolder),
		RestoreIpfsFolder:  generateRestoreIpfsFolder(storeIpfsFolder),
		UploadIpfsFolder:   generateUploadIpfsFolder(storeIpfsFolder),
	}
	return config
}

// Lee las variables de entorno que establecen la configuración de almacenamiento and clean the store
func LoadConfig() {
	inputDataFolder := os.Getenv("INPUT_DATA_FOLDER")
	storeIpfsFolder := os.Getenv("STORE_IPFS_FOLDER")
	config := NewStorageConfig(inputDataFolder, storeIpfsFolder)
	config.LogConfig()
	LocalStorageConfig = *config
}

// Hace un log por debug de la configuración
func (config *StorageConfig) LogConfig() {
	common.Log.Debugf("InputDataFolder: %v | StoreIpfsFolder: %v | DownloadIpfsFolder: %v | RestoreIpfsFolder: %v | UploadIpfsFolder: %v",
		config.InputDataFolder,
		config.StoreIpfsFolder,
		config.DownloadIpfsFolder,
		config.RestoreIpfsFolder,
		config.UploadIpfsFolder,
	)
}

// Genera el directorio donde se guardan las descargas
func generateDownloadIpfsFolder(storeIpfsFolder string) string {
	return storeIpfsFolder + "/" + DOWLOAD_SUB_DIRECTORY
}

// Genera el directorio donde se guardan los archivos restaurados
func generateRestoreIpfsFolder(storeIpfsFolder string) string {
	return storeIpfsFolder + "/" + RESTORE_SUB_DIRECTORY
}

// Genera el directorio donde se guardan los archivos subidos desde el exterior de la red
func generateUploadIpfsFolder(storeIpfsFolder string) string {
	return storeIpfsFolder + "/" + UPLOAD_SUB_DIRECTORY
}
