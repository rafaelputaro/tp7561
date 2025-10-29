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

// Contiene las métricas del módulo cliente
type Metrics struct {
	uploadedFileCount prometheus.Counter
}

// Aloja un buffer con las métricas del sistema las cuales puede ser recuperadas.
// Se interactúa con el mismo a través de una API-REST
type MetricsService struct {
	config  common_metrics.MetricsConfig
	metrics Metrics
	reg     *prometheus.Registry
}

// Retorna un servidor de métricas de cliente listo para ser utilizado
func NewMetricsServer(namespace string) *MetricsService {
	config := common_metrics.LoadMetricsConfig(namespace)
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

/*

scrape_configs:
  - job_name: 'nodos_mensajes'
    static_configs:
      - targets: ['nodo1:9100', 'nodo2:9100', 'nodo3:9100']


*/

/*
connected_nodes{node_id="nodo1", direction="outbound"} 1
connected_nodes{node_id="nodo2", direction="inbound"} 1

*/

/*

scrape_configs:
  - job_name: 'nodos'
    static_configs:
      - targets: ['nodo1:8080', 'nodo2:8080', 'nodo3:8080']

En Grafana, crea un panel con una consulta PromQL como:

connected_nodes

Usa un gráfico de tipo "Bar gauge" o "Table" para mostrar qué nodos tienen valor 1 (conectados). También puedes usar una consulta como:

up{job="nodos"}

para ver qué nodos son alcanzables por Prometheus (estado up).

*/
