package client_metrics

import (
	"net/http"
	"tp/common"
	common_metrics "tp/common/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const MSG_CLOSING_SERVICE = "closing service: %v"
const MSG_SERVICE_STARTED = "service started: %v"

var MetricsServiceInstance = NewMetricsServer()

// Contiene las métricas del módulo cliente
type Metrics struct {
	uploadedFileCount prometheus.Counter
}

// Aloja las métricas del sistema las cuales puede ser recuperadas.
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

// Retorna una nueva instancia de métrica de cliente lista para ser utilizada
func newMetrics(namespace string, reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		uploadedFileCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "uploaded_file_count",
			Help:      "Number of file uploaded byt the client",
		}),
	}
	reg.MustRegister(m.uploadedFileCount)
	return m
}

// Se encarga de antender las solicitudes de prometheus
func (service *MetricsService) Serve() {
	common.Log.Infof(MSG_SERVICE_STARTED, service.config.UrlPrometheus)
	promHandler := promhttp.HandlerFor(service.reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", promHandler)
	http.ListenAndServe(service.config.UrlPrometheus, nil)
	common.Log.Infof(MSG_CLOSING_SERVICE, service.config.UrlPrometheus)
}

// Incrementa en uno la cantidad de archivos subidos desde este módulo al sistema
func (service *MetricsService) IncUploadedFileCount() {
	service.metrics.uploadedFileCount.Inc()
}
