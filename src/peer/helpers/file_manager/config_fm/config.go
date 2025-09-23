package config_fm

import (
	"os"
	"tp/common"
)

type StorageConfig struct {
	InputDataFolder string
	StoreIpfsFolder string
}

var LocalStorageConfig StorageConfig

// Retorna una nueva instancia de la configuración
func NewStorageConfig(inputDataFolder string, storeIpfsFolder string) *StorageConfig {
	config := &StorageConfig{
		InputDataFolder: inputDataFolder,
		StoreIpfsFolder: storeIpfsFolder,
	}
	return config
}

// Lee las variables de entorno que establecen la configuración de almacenamiento
func LoadConfig() {
	inputDataFolder := os.Getenv("INPUT_DATA_FOLDER")
	storeIpfsFolder := os.Getenv("STORE_IPFS_FOLDER")
	config := NewStorageConfig(inputDataFolder, storeIpfsFolder)
	config.LogConfig()
	LocalStorageConfig = *config
}

// Hace un log por debug de la configuración
func (config *StorageConfig) LogConfig() {
	common.Log.Debugf("InputDataFolder: %v | StoreIpfsFolder: %v",
		config.InputDataFolder,
		config.StoreIpfsFolder,
	)
}
