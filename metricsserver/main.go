package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var pingCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "ping_request_count",
		Help: "No of request handled by Ping handler",
	},
)

var cpuUsage = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "cpu",
		Help: "CPU utilization",
	},
)

var reqPerSecond = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "reqpersecond",
		Help: "Requests per second",
	},
)

var scaleToZero = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "scaletozero",
		Help: "0 if wanting to scale to zero otherwise 1 to not scale to zero",
	},
)

func ping(w http.ResponseWriter, req *http.Request) {
	pingCounter.Inc()
	fmt.Fprintf(w, "pong")
}

func set(w http.ResponseWriter, req *http.Request) {
	metric := req.URL.Query().Get("metric")
	value := req.URL.Query().Get("value")

	if metric == "" || value == "" {
		fmt.Fprintf(w, "metric and value both need to be set\n")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// string to int
	valueInt, err := strconv.Atoi(value)
	if err != nil {
		fmt.Fprintf(w, "value must be an integer\n")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch metric {
	case "scaletozero":
		scaleToZero.Set(float64(valueInt))
	case "cpu":
		cpuUsage.Set(float64(valueInt))
	case "reqpersecond":
		reqPerSecond.Set(float64(valueInt))
	default:
		fmt.Fprintf(w, "metric query param must be 'cpu' or 'scaletozero' or 'reqpersecond'\n")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "metric %s set at %d\n", metric, valueInt)
}

func main() {
	prometheus.MustRegister(pingCounter)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(reqPerSecond)
	prometheus.MustRegister(scaleToZero)

	fmt.Println("Endpoints:")
	fmt.Println("/set (query param 'name' and 'value') -> set a metric at value given")
	fmt.Println("/metrics -> dump metrics")
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/set", set)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("Listening on :8090...")
	_ = http.ListenAndServe(":8090", nil)
}
