package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

type metrics struct {
	devices prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "metrics",
			Name:      "connected_devices",
			Help:      "Number....",
		}),
	}
	reg.MustRegister(m.devices)
	return m
}

var dvs []Device

func init() {
	dvs = []Device{
		{1, "5F-33-CC-IF-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2-1.6"},
	}
}

func main() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	m.devices.Set(float64(len(dvs)))

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	http.Handle("/metrics", promHandler)
	http.HandleFunc("/devices", getDevices)
	http.ListenAndServe(":9392", nil)
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(dvs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
