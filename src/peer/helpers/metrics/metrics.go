package peer_metrics

import "github.com/prometheus/client_golang/prometheus"

const SAVE_METRIC = "Save metric: %v | value: %v"

// Contiene las métricas del módulo cliente
type Metrics struct {
	lastFileReturnedNumberMetric *LastFileReturnedNumberMetric
	contactMetrics               *ContacMetrics
	contactCount                 *ContactCountMetric
}

// Retorna una nueva instancia de métrica de cliente lista para ser utilizada
func newMetrics(namespace string, reg prometheus.Registerer) *Metrics {
	return &Metrics{
		lastFileReturnedNumberMetric: newLastFileReturnedNumberMetric(namespace, reg),
		contactMetrics:               newContactMetrics(namespace, reg),
		contactCount:                 newContacCountMetric(namespace, reg, parseUrlToId(namespace)),
	}
}
