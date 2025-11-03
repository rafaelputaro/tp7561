package peer_metrics

import (
	"fmt"
	"tp/common"
	"tp/common/contact"

	"github.com/prometheus/client_golang/prometheus"
)

const CONTACT_METRICS_NAME = "contact_list"
const CONTACT_METRICS_HELP = "Peer contacts"
const ADD_CONTACT_METRIC = "add contact | source: %v | target: %v"
const REMOVE_CONTACT_METRIC = "remove contact | source: %v | target: %v"

// Representa la métrica del último archivo o bloque retornado por el par actual
type ContacMetrics struct {
	contacsVec *prometheus.CounterVec
}

// Retorna una nueva instancia de esta métrica
func newContactMetrics(namespace string, reg prometheus.Registerer) *ContacMetrics {
	contactVec := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      CONTACT_METRICS_NAME,
		Help:      CONTACT_METRICS_HELP,
	},
		[]string{"source", "target"},
	)
	m := &ContacMetrics{
		contacsVec: contactVec,
	}
	reg.MustRegister(m.contacsVec)
	return m
}

// Agregar un contacto
func (metric *ContacMetrics) addContact(sourceName string, target contact.Contact) {
	targetName := parseContact(target)
	value := fmt.Sprintf(ADD_CONTACT_METRIC, sourceName, targetName)
	common.Log.Debugf(SAVE_METRIC, CONTACT_METRICS_NAME, value)
	metric.contacsVec.WithLabelValues(sourceName, targetName).Inc()
}

// Remover un contacto
func (metric *ContacMetrics) removeContact(sourceName string, target contact.Contact) {
	targetName := parseContact(target)
	value := fmt.Sprintf(REMOVE_CONTACT_METRIC, sourceName, targetName)
	common.Log.Debugf(SAVE_METRIC, CONTACT_METRICS_NAME, value)
	metric.contacsVec.DeleteLabelValues(sourceName, targetName)
}
