package helpers

import (
	"os"
	"strconv"
)

const DEFAULT_ENTRIES_PER_K_BUCKET = 10
const MSG_ERROR_ON_LOAD_ENTRIES_PER_K_BUCKET = "Error on load entries per k bucket"
const EMPTY_URL = ""

// Representa la configuración del par
type PeerConfig struct {
	Name              string
	Url               string
	EntriesPerKBucket int
}

// Retorna una nueva instancia de la configuración
func NewNodeConfig(name string, url string, entriesPerKBucket int) *PeerConfig {
	config := &PeerConfig{
		Name:              name,
		Url:               url,
		EntriesPerKBucket: entriesPerKBucket,
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
	entries_per_k_bucket, err := strconv.Atoi(entries_per_k_bucket_s)
	if err != nil {
		entries_per_k_bucket = DEFAULT_ENTRIES_PER_K_BUCKET
		Log.Debugf(MSG_ERROR_ON_LOAD_ENTRIES_PER_K_BUCKET)
	}
	var config = NewNodeConfig(name, GenerateURL(host, port), entries_per_k_bucket)
	config.LogConfig()
	return config
}

// Hace un log por debug de la configuración
func (config *PeerConfig) LogConfig() {
	Log.Debugf("Name: %v | Url: %v | EntriesPerKBucket: %v",
		config.Name,
		config.Url,
		config.EntriesPerKBucket,
	)
}
