package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var gauges map[string]prometheus.Gauge

func set(w http.ResponseWriter, req *http.Request) {
	if len(req.URL.Query()) == 0 {
		fmt.Fprintf(w, "metrics should be provided as query parameters (e.g. /set?cpu=50)\n")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	for name, value := range req.URL.Query() {
		if name == "" || len(value) != 1 || value[0] == "" {
			fmt.Fprintf(w, "metrics name and value should both be set\n")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		// string to int
		value, err := strconv.Atoi(value[0])
		if err != nil {
			fmt.Fprintf(w, "value must be an integer\n")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = setMetric(name, float64(value))
		if err != nil {
			fmt.Fprintf(w, "could not set metric %s: %s", name, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "metric %s set at %d\n", name, value)
	}
}

func setMetric(name string, value float64) error {
	metric, ok := gauges[name]
	if !ok {
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: name,
				Help: fmt.Sprintf("Dynamically created %s", name),
			},
		)
		err := prometheus.Register(metric)
		if err != nil {
			return err
		}

		gauges[name] = metric
	}

	metric.Set(value)
	return nil
}

func unset(w http.ResponseWriter, req *http.Request) {
	if len(req.URL.Query()) == 0 {
		fmt.Fprintf(w, "metrics should be provided as query parameters (e.g. /unset?cpu)\n")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	for name := range req.URL.Query() {
		if name == "" {
			fmt.Fprintf(w, "metrics name should be set\n")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		unsetMetric(name)
		fmt.Fprintf(w, "unset metric %s\n", name)
	}
}

func unsetMetric(name string) {
	metric, ok := gauges[name]
	if ok {
		delete(gauges, name)
		prometheus.Unregister(metric)
	}
}

func main() {
	gauges = map[string]prometheus.Gauge{}

	fmt.Println("Endpoints:")
	fmt.Println("/set -> register one or more metrics. Example: /set?cpu=55&memory=1200")
	fmt.Println("/unset -> unregister one or more metrics. Example: /unset?cpu&memory")
	fmt.Println("/metrics -> dump metrics")
	http.HandleFunc("/set", set)
	http.HandleFunc("/unset", unset)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Listening on :8090...")
	_ = http.ListenAndServe(":8090", nil)
}
