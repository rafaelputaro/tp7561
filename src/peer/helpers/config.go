package helpers

import (
	"fmt"
	"os"
	"strconv"
	"tp/common"
	"tp/common/keys"
)

const DEFAULT_ENTRIES_PER_K_BUCKET = 10
const DEFAULT_SEARCH_WORKERS = 10
const DEFAULT_NUMBER_OF_PAIRS = 15
const MSG_ERROR_ON_LOAD_ENTRIES_PER_K_BUCKET = "Error on load entries per k bucket"
const MSG_ERROR_ON_LOAD_SEARCH_WORKERS = "Error on load search workers"
const MSG_ERROR_ON_LOAD_NUMBER_OF_PAIRS = "Error on load number of pairs"
const EMPTY_URL = ""
const RECOVERED_FOLDER = "recovered"

// Representa la configuración del par
type PeerConfig struct {
	Id                []byte
	Name              string
	Url               string
	Port              string
	EntriesPerKBucket int
	SearchWorkers     int
	NumberOfPairs     int
}

// Retorna una nueva instancia de la configuración
func NewNodeConfig(name string, url string, port string, entriesPerKBucket int, searchWorkers int, numberOfPairs int) *PeerConfig {
	config := &PeerConfig{
		Id:                keys.GetKey(name),
		Name:              name,
		Url:               url,
		Port:              port,
		EntriesPerKBucket: entriesPerKBucket,
		SearchWorkers:     searchWorkers,
		NumberOfPairs:     numberOfPairs,
	}
	return config
}

func GenerateURL(host string, port string) string {
	return host + ":" + port
}

// Lee las variables de entorno que establecen la configuración del par
func LoadConfig() *PeerConfig {
	name := os.Getenv("PEER_NAME")
	port := os.Getenv("PEER_PORT")
	host := os.Getenv("PEER_HOST")
	entries_per_k_bucket_s := os.Getenv("ENTRIES_PER_K_BUCKET")
	search_workers_s := os.Getenv("SEARCH_WORKERS")
	number_of_pairs_s := os.Getenv("NUMBER_OF_PAIRS")
	entries_per_k_bucket, err := strconv.Atoi(entries_per_k_bucket_s)
	if err != nil {
		entries_per_k_bucket = DEFAULT_ENTRIES_PER_K_BUCKET
		common.Log.Debugf(MSG_ERROR_ON_LOAD_ENTRIES_PER_K_BUCKET)
	}
	search_workers, err := strconv.Atoi(search_workers_s)
	if err != nil {
		search_workers = DEFAULT_SEARCH_WORKERS
		common.Log.Debugf(MSG_ERROR_ON_LOAD_SEARCH_WORKERS)
	}
	number_of_pairs, err := strconv.Atoi(number_of_pairs_s)
	if err != nil {
		search_workers = DEFAULT_NUMBER_OF_PAIRS
		common.Log.Debugf(MSG_ERROR_ON_LOAD_NUMBER_OF_PAIRS)
	}
	var config = NewNodeConfig(name, GenerateURL(host, port), port, entries_per_k_bucket, search_workers, number_of_pairs)
	config.LogConfig()
	return config
}

// Hace un log por debug de la configuración
func (config *PeerConfig) LogConfig() {
	common.Log.Debugf("Name: %v | Url: %v | Id: %v | EntriesPerKBucket: %v | SearchWorkers: %v | NumberOfPairs: %v",
		config.Name,
		config.Url,
		fmt.Sprintf("%v", keys.KeyToLogFormatString(config.Id)),
		config.EntriesPerKBucket,
		config.SearchWorkers,
		config.NumberOfPairs,
	)
}
