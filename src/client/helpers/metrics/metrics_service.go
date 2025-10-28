package client_metrics

import (
	"net/http"
	"tp/common"
	common_metrics "tp/common/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const MSG_ERROR_CLOSING_SERVICE = "error closing service: %v"
const MSG_SERVICE_ERROR = "service error: %v"
const MSG_CLOSING_SERVICE = "closing service"
const MSG_SERVICE_STARTED = "service started: %v"

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

func NewMetricsServer() *MetricsService {
	config := common_metrics.LoadMetricsConfig()
	reg := prometheus.NewRegistry()
	server := MetricsService{
		config:  *config,
		metrics: *NewMetrics(reg),
		reg:     reg,
	}
	return &server
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		uploadedFileCount: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "client-1",
			Name:      "uploaded_file_count",
			Help:      "Number of file uploaded byt the client",
		}),
	}
	reg.MustRegister(m.uploadedFileCount)
	return m
}

func (service *MetricsService) Serve() {
	common.Log.Infof(MSG_SERVICE_STARTED, service.config.Url)
	promHandler := promhttp.HandlerFor(service.reg, promhttp.HandlerOpts{})

	http.Handle("/metrics", promHandler)
	http.ListenAndServe(service.config.UrlPrometheus, nil)
	common.Log.Infof("End service: %v", service.config.Url)
}

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
