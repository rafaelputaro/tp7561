package common_metrics

import (
	"os"
	"tp/common"
	"tp/common/communication/url"
)

// Contiene la configuración de las métricas
type MetricsConfig struct {
	Url string
}

// Retorna la configuración de las métricas
func NewMetricsConfig(host string, port string) *MetricsConfig {
	return &MetricsConfig{
		Url: url.GenerateURL(host, port),
	}
}

// Lee las variables de entorno que establecen la configuración las métricas
func LoadMetricsConfig() *MetricsConfig {
	host := os.Getenv("METRICS_HOST")
	port := os.Getenv("METRICS_PORT")
	var config = NewMetricsConfig(host, port)
	config.LogConfig()
	return config
}

// Hace un log por debug de la configuración
func (config *MetricsConfig) LogConfig() {
	common.Log.Debugf("Url: %v", config.Url)
}
