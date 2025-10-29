package common_metrics

import (
	"os"
	"tp/common"
	"tp/common/communication/url"
)

// Contiene la configuración de las métricas
type MetricsConfig struct {
	Namespace     string
	Url           string
	UrlPrometheus string
}

// Retorna la configuración de las métricas
func NewMetricsConfig(namespace string, host string, port string) *MetricsConfig {
	return &MetricsConfig{
		Namespace:     namespace,
		Url:           url.GenerateURL(host, port),
		UrlPrometheus: generateUrlPrometheus(port),
	}
}

// Lee las variables de entorno que establecen la configuración las métricas
func LoadMetricsConfig(namespace string) *MetricsConfig {
	host := os.Getenv("METRICS_HOST")
	port := os.Getenv("METRICS_PORT")
	var config = NewMetricsConfig(namespace, host, port)
	config.LogConfig()
	return config
}

// Hace un log por debug de la configuración
func (config *MetricsConfig) LogConfig() {
	common.Log.Debugf("Namespace: %v | Url: %v | UrlPrometheus: %v", config.Namespace, config.Url, config.UrlPrometheus)
}

// Retorna la url a ser utilizada para atender los request de prometheus
func generateUrlPrometheus(port string) string {
	return ":" + port
}
