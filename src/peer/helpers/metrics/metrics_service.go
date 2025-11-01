package peer_metrics

import (
	"net/http"
	"tp/common"
	"tp/common/contact"
	common_metrics "tp/common/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const MSG_CLOSING_SERVICE = "closing service: %v"
const MSG_SERVICE_STARTED = "service started: %v"

var metricsServiceInstance = NewMetricsServer()

// Aloja las métricas del sistema las cuales puede ser recuperadas por medio de Prometheus.
type MetricsService struct {
	config  common_metrics.MetricsConfig
	metrics Metrics
	reg     *prometheus.Registry
}

// Retorna un servidor de métricas de cliente listo para ser utilizado
func NewMetricsServer() *MetricsService {
	config := common_metrics.LoadMetricsConfig()
	reg := prometheus.NewRegistry()
	server := MetricsService{
		config:  *config,
		metrics: *newMetrics(config.Namespace, reg),
		reg:     reg,
	}
	return &server
}

// Se encarga de antender las solicitudes de prometheus
func Serve() {
	common.Log.Infof(MSG_SERVICE_STARTED, metricsServiceInstance.config.UrlPrometheus)
	promHandler := promhttp.HandlerFor(metricsServiceInstance.reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", promHandler)
	http.ListenAndServe(metricsServiceInstance.config.UrlPrometheus, nil)
	common.Log.Infof(MSG_CLOSING_SERVICE, metricsServiceInstance.config.UrlPrometheus)
}

// Incrementa en uno la cantidad de archivos subidos desde este módulo al sistema
func SetLastFileReturnedNumber(fileName string) {
	metricsServiceInstance.metrics.lastFileReturnedNumberMetric.setLastFileReturnedNumber(fileName)
}

// Agregar un contacto
func AddContact(sourceName string, target contact.Contact) {
	metricsServiceInstance.metrics.contactMetrics.addContact(sourceName, target)
	metricsServiceInstance.metrics.contactCount.incCount()
}

// Remover un contacto
func RemoveContact(sourceName string, target contact.Contact) {
	metricsServiceInstance.metrics.contactMetrics.removeContact(sourceName, target)
	metricsServiceInstance.metrics.contactCount.descCount()
}
