package helpers

import (
	"fmt"
	"os"
	"strconv"
	"tp/common"
)

const DEFAULT_ENTRIES_PER_K_BUCKET = 10
const DEFAULT_SEARCH_WORKERS = 10
const MSG_ERROR_ON_LOAD_ENTRIES_PER_K_BUCKET = "Error on load entries per k bucket"
const MSG_ERROR_ON_LOAD_SEARCH_WORKERS = "Error on load search workers"
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
}

// Retorna una nueva instancia de la configuración
func NewNodeConfig(name string, url string, port string, entriesPerKBucket int, searchWorkers int) *PeerConfig {
	config := &PeerConfig{
		Id:                GetKey(name),
		Name:              name,
		Url:               url,
		Port:              port,
		EntriesPerKBucket: entriesPerKBucket,
		SearchWorkers:     searchWorkers,
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
	var config = NewNodeConfig(name, GenerateURL(host, port), port, entries_per_k_bucket, search_workers)
	config.LogConfig()
	return config
}

// Hace un log por debug de la configuración
func (config *PeerConfig) LogConfig() {
	common.Log.Debugf("Name: %v | Url: %v | Id: %v | EntriesPerKBucket: %v | SearchWorkers: %v",
		config.Name,
		config.Url,
		fmt.Sprintf("%v", KeyToLogFormatString(config.Id)),
		config.EntriesPerKBucket,
		config.SearchWorkers,
	)
}
