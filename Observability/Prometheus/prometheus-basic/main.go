package main

import (
	"fmt"
	"html"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type metrics struct {
	pingCounter           prometheus.Counter
	cpuTemp               prometheus.Gauge
	hdFailures            *prometheus.CounterVec
	responseTimeHistogram *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		pingCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "ping_request_count",
				Help: "No of request handled by Ping handler",
			},
		),
		cpuTemp: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "cpu_temperature_celsius",
				Help: "Current temperature of the CPU.",
			},
		),
		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hd_errors_total",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		),
		responseTimeHistogram: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: "namespace",
				Name:      "http_server_request_duration_seconds",
				Help:      "Histogram of response time for handler in seconds",
				Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"route", "method", "status_code"},
		),
	}
	reg.MustRegister(m.pingCounter)
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.hdFailures)
	reg.MustRegister(m.responseTimeHistogram)
	return m
}

// statusRecorder to record the status code from the ResponseWriter
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(statusCode int) {
	rec.statusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func measureResponseDuration(next http.Handler, m metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := statusRecorder{w, 200}

		sleepTime := rand.Intn(3000)
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
		log.Println("Location : ", r.URL.Path, "Time", sleepTime)

		next.ServeHTTP(&rec, r)

		duration := time.Since(start)
		statusCode := strconv.Itoa(rec.statusCode)
		route := r.RequestURI
		m.responseTimeHistogram.WithLabelValues(route, r.Method, statusCode).Observe(duration.Seconds())
	})
}

func ping(w http.ResponseWriter, req *http.Request, m *metrics) {
	m.pingCounter.Inc()
	m.pingCounter.Add(rand.Float64())

	fmt.Fprintf(w, "pong")
}

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)
	// Set values for the new created metrics.
	m.cpuTemp.Set(65.3)
	m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		ping(w, r, m)

		fmt.Fprintf(w, "Ping")
		log.Println("Ping Log")
	})

	http.Handle("/", measureResponseDuration(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.RequestURI)
	}), *m))
	http.Handle("/bar", measureResponseDuration(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.RequestURI)
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		tm := time.Now().String()
		w.Write([]byte("The time is: " + tm))
	}), *m))

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.RequestURI)
		log.Println(r.RequestURI)
		tm := time.Now().String()
		w.Write([]byte("The time is: " + tm))
	})

	log.Fatal(http.ListenAndServe(":8090", nil))
}
