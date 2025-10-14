package helpers

import (
	"os"
	"strconv"
	"tp/common"
	"tp/common/communication/url"
)

const DEFAULT_NUMBER_OF_PAIRS = 15
const DEFAULT_NUMBER_OF_CLIENTS = 15
const MSG_ERROR_ON_LOAD_NUMBER_OF_PAIRS = "Error on load number of pairs"
const MSG_ERROR_ON_LOAD_NUMBER_OF_CLIENTS = "Error on load number of clients"
const EMPTY_URL = ""
const RECOVERED_FOLDER = "recovered"
const DOWLOAD_SUB_DIRECTORY = "down"
const RESTORE_SUB_DIRECTORY = "restore"

// Representa la configuración del client
type Config struct {
	Name            string
	Url             string
	Port            string
	NumberOfClients int
	NumberOfPairs   int
	InputDataFolder string
	StoreFolder     string
	DownloadFolder  string
	RestoreFolder   string
}

// Retorna una nueva instancia de la configuración
func NewConfig(name string, url string, port string, numberOfClients int, numberOfPairs int, inputDataFolder string, storeFolder string) *Config {
	config := &Config{
		Name:            name,
		Url:             url,
		Port:            port,
		NumberOfClients: numberOfClients,
		NumberOfPairs:   numberOfPairs,
		InputDataFolder: inputDataFolder,
		StoreFolder:     storeFolder,
		DownloadFolder:  generateDownloadFolder(storeFolder),
		RestoreFolder:   generateRestoreFolder(storeFolder),
	}
	return config
}

// Lee las variables de entorno que establecen la configuración del client
func LoadConfig() *Config {
	name := os.Getenv("CLIENT_NAME")
	port := os.Getenv("CLIENT_PORT")
	host := os.Getenv("CLIENT_HOST")
	number_of_clients_s := os.Getenv("NUMBER_OF_CLIENTS")
	number_of_pairs_s := os.Getenv("NUMBER_OF_PAIRS")
	// conversión clients
	number_of_clients, err := strconv.Atoi(number_of_clients_s)
	if err != nil {
		number_of_clients = DEFAULT_NUMBER_OF_CLIENTS
		common.Log.Debugf(MSG_ERROR_ON_LOAD_NUMBER_OF_CLIENTS)
	}
	// conversión pares
	number_of_pairs, err := strconv.Atoi(number_of_pairs_s)
	if err != nil {
		number_of_pairs = DEFAULT_NUMBER_OF_PAIRS
		common.Log.Debugf(MSG_ERROR_ON_LOAD_NUMBER_OF_PAIRS)
	}
	// leer directorios
	inputDataFolder := os.Getenv("INPUT_DATA_FOLDER")
	storeFolder := os.Getenv("STORE_FOLDER")
	var config = NewConfig(name, url.GenerateURL(host, port), port, number_of_clients, number_of_pairs, inputDataFolder, storeFolder)
	config.LogConfig()
	return config
}

// Hace un log por debug de la configuración
func (config *Config) LogConfig() {
	common.Log.Debugf("Name: %v | Url: %v | NumberOfClients: %v | NumberOfPairs: %v | InputDataFolder: %v | StoreFolder: %v | DownloadFolder: %v | RestoreFolder: %v ",
		config.Name,
		config.Url,
		config.NumberOfClients,
		config.NumberOfPairs,
		config.InputDataFolder,
		config.StoreFolder,
		config.DownloadFolder,
		config.RestoreFolder,
	)
}

// Genera el directorio donde se guardan las descargas
func generateDownloadFolder(storeFolder string) string {
	return storeFolder + "/" + DOWLOAD_SUB_DIRECTORY
}

// Genera el directorio donde se guardan los archivos restaurados
func generateRestoreFolder(storeFolder string) string {
	return storeFolder + "/" + RESTORE_SUB_DIRECTORY
}
