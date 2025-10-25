package main

import (
	"os"
	"tp/common"
	"tp/common/communication/url"
)

// Contiene la configuración del servidor
type MetricsServerConfig struct {
	Url string
}

// Retorna la configuración del servidor
func NewMetricsServerConfig(host string, port string) *MetricsServerConfig {
	return &MetricsServerConfig{
		Url: url.GenerateURL(host, port),
	}
}

// Lee las variables de entorno que establecen la configuración del servidor
func LoadMetricsServerConfig() *MetricsServerConfig {
	host := os.Getenv("METRICS_HOST")
	port := os.Getenv("METRICS_PORT")
	var config = NewMetricsServerConfig(host, port)
	config.LogConfig()
	return config
}

// Hace un log por debug de la configuración
func (config *MetricsServerConfig) LogConfig() {
	common.Log.Debugf("Url: %v", config.Url)
}
