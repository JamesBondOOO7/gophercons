package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// https://youtu.be/1V7eJ0jN8-E?si=eyrhKxjMZfoa4cHR
// https://www.youtube.com/watch?v=WUBjlJzI2a0&t=135s
// https://www.digitalocean.com/community/tutorials/understanding-init-in-go

type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

type metrics struct {
	devices prometheus.Gauge     // representing a single value for a metric
	info    *prometheus.GaugeVec // exposing arbitrary # of key-value pairs
}

var dvs []Device
var version string

func init() {
	version = "2.10.5"
	dvs = []Device{
		{1, "5F-33-CC-1F-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2.1.6"},
	}
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

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "connected_devices",
			Help:      "Number of currently connected devices.",
		}),
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "myapp",
			Name:      "info",
			Help:      "Info about the My App environment.",
		},
			[]string{"version"}),
	}
	reg.MustRegister(m.devices, m.info)
	return m
}

type registerDevicesHandler struct {
	metrics *metrics
}

func (rdh *registerDevicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getDevices(w, r)
	case "POST":
		createDevice(w, r, rdh.metrics)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createDevice(w http.ResponseWriter, r *http.Request, m *metrics) {
	var dv Device
	err := json.NewDecoder(r.Body).Decode(&dv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dvs = append(dvs, dv)
	m.devices.Set(float64(len(dvs)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Device Created!"))
}

func main() {
	reg := prometheus.NewRegistry()
	m := NewMetrics(reg)
	m.devices.Set(float64(len(dvs)))
	m.info.With(prometheus.Labels{"version": version}).Set(1)

	dMux := http.NewServeMux()
	rdh := &registerDevicesHandler{m}
	dMux.Handle("/devices", rdh)

	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	// Exposing business logic under :8080 port
	go func() {
		log.Fatal(http.ListenAndServe(":8080", dMux))
	}()

	// For safety & control, exposing metrics under :8081 port
	go func() {
		log.Fatal(http.ListenAndServe(":8081", pMux))
	}()

	select {}
}
