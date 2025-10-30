package peer_metrics

import (
	"net/http"
	"tp/common"
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

/*

// Crear un GaugeFunc que se evalúa en cada scrape
	gauge := prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "mi_gauge_personalizado",
			Help: "Gauge cuyo valor se modifica al momento del scrape",
		},
		func() float64 {
			// Aquí puedes modificar el valor antes de devolverlo
			// Ejemplo: decrementar en 10 cada vez que se lee
			externalGaugeValue -= 10
			return externalGaugeValue
		},
	)

*/

/*

const MSG_CLOSING_SERVICE = "closing service: %v"
const MSG_SERVICE_STARTED = "service started: %v"

var metricsServiceInstance = NewMetricsServer()

// Contiene las métricas del módulo cliente
type Metrics struct {
	lastFileReturnedNumber prometheus.Gauge
}

// Aloja las métricas del sistema las cuales puede ser recuperadas.
type MetricsService struct {
	config       common_metrics.MetricsConfig
	metrics      Metrics
	reg          *prometheus.Registry
	mutexMetrics sync.Mutex
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
		lastFileReturnedNumber: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "last_file_returned_number",
			Help:      "Number associated with the last returned file or block",
		}),
	}
	reg.MustRegister(m.lastFileReturnedNumber)


	return m
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
	metricsServiceInstance.mutexMetrics.Lock()
	defer metricsServiceInstance.mutexMetrics.Unlock()
	common.Log.Debugf("Save metric: %v | value: %v", metricsServiceInstance.config.Namespace, fileName)
	metricsServiceInstance.metrics.lastFileReturnedNumber.Set(float64(parseFileNumber(fileName)))
}

// Retorna el segundo número que aparece en el nombre de un archivo
func parseFileNumber(fileName string) int {
	patron := `(\d+)`
	ren := regexp.MustCompile(patron)
	found := ren.FindAllString(fileName, -1)
	common.Log.Debugf("%v", found)
	if len(found) > 1 {
		converted, err := strconv.Atoi(found[1])
		if err == nil {
			return converted
		}
	}
	return -1
}

*/
