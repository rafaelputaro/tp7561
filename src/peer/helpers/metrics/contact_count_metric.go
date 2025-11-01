package peer_metrics

import (
	"fmt"
	"tp/common"

	"github.com/prometheus/client_golang/prometheus"
)

const CONTACT_COUNT_METRIC_NAME = "contact_count"
const CONTACT_COUNT_METRIC_HELP = "Contact count"
const INC_CONTACT_COUNT_METRIC = "inc contact count | id: %v"
const DEC_CONTACT_COUNT_METRIC = "dec contact count | id: %v"
const MSG_ERROR_ON_INC_COUNT = "error on inc contact count: %v"
const MSG_ERROR_ON_DEC_COUNT = "error on dec contact count: %v"

// Representa la métrica del último archivo o bloque retornado por el par actual
type ContactCountMetric struct {
	contactCount *prometheus.CounterVec
	id           string
}

// Retorna una nueva instancia de esta métrica
func newContacCountMetric(namespace string, reg prometheus.Registerer, id string) *ContactCountMetric {
	count := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      CONTACT_COUNT_METRIC_NAME,
		Help:      CONTACT_COUNT_METRIC_HELP,
	}, []string{"id"})
	m := &ContactCountMetric{
		id:           id,
		contactCount: count,
	}

	reg.MustRegister(m.contactCount)
	m.contactCount.WithLabelValues(id).Add(0)
	return m
}

func (metric *ContactCountMetric) incCount() {
	id := metric.id
	value := fmt.Sprintf(INC_CONTACT_COUNT_METRIC, id)
	common.Log.Debugf(SAVE_METRIC, CONTACT_COUNT_METRIC_NAME, value)
	metric.contactCount.WithLabelValues(id).Inc()
}

func (metric *ContactCountMetric) descCount() {
	id := metric.id
	value := fmt.Sprintf(DEC_CONTACT_COUNT_METRIC, id)
	common.Log.Debugf(SAVE_METRIC, CONTACT_COUNT_METRIC_NAME, value)
	metric.contactCount.WithLabelValues(id).Desc()
}
