package peer_metrics

import (
	"sync"
	"tp/common"
	common_metrics "tp/common/metrics"

	"github.com/prometheus/client_golang/prometheus"
)

const INVALID_VALUE = -1
const LAST_FILE_RETURNED_NUMBER_NAME = "last_file_returned_number"
const LAST_FILE_RETURNED_NUMBER_HELP = "Number associated with the last returned file or block"

// Representa la métrica del último archivo o bloque retornado por el par actual
type LastFileReturnedNumberMetric struct {
	lastFileReturnedNumberFunc prometheus.GaugeFunc
	lastFileReturnedNumber     float64
	mutex                      sync.Mutex
}

// Retorna una nueva instancia de esta métrica
func newLastFileReturnedNumberMetric(namespace string, reg prometheus.Registerer) *LastFileReturnedNumberMetric {
	m := &LastFileReturnedNumberMetric{
		lastFileReturnedNumber: INVALID_VALUE,
	}
	m.lastFileReturnedNumberFunc = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      LAST_FILE_RETURNED_NUMBER_NAME,
		Help:      LAST_FILE_RETURNED_NUMBER_HELP,
	}, func() float64 {
		m.mutex.Lock()
		defer m.mutex.Unlock()
		last := m.lastFileReturnedNumber
		m.lastFileReturnedNumber = INVALID_VALUE
		return last
	})
	reg.MustRegister(m.lastFileReturnedNumberFunc)
	return m
}

// Incrementa en uno la cantidad de archivos subidos desde este módulo al sistema
func (metric *LastFileReturnedNumberMetric) setLastFileReturnedNumber(fileName string) {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()
	parsed := common_metrics.ParseFileNumber(fileName)
	common.Log.Debugf(SAVE_METRIC, LAST_FILE_RETURNED_NUMBER_NAME, parsed)
	metric.lastFileReturnedNumber = parsed
}
