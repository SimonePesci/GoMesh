package proxy

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {

	// Counter for the total number of requests (by service and status code)
	RequestsTotal *prometheus.CounterVec

	// Tracks request latency distribution
	RequestDuration *prometheus.HistogramVec

	// Tracks how many requests are currently being handled by the proxy
	RequestsInFlight prometheus.Gauge

	// Tracks the number of errors (by type)
	ErrorsTotal *prometheus.CounterVec
}

func NewMetrics() *Metrics {

	metrics := &Metrics{
		// Counter for the total number of requests 
		// Labeled by service (which backend) and status code
		// Using promauto to automatically register with the default registry
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gomesh_requests_total",
				Help: "Total number of requests received by the proxy",
			},
			[]string{"service", "status"},
		),

		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "gomesh_request_duration_seconds",
				Help: "Requests duration in seconds",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
			},
			[]string{"service"},
		),

		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "gomesh_requests_in_flight",
				Help: "Number of requests currently being handled by the proxy",
			},
		),

		ErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "gomesh_errors_total",
				Help: "Total number of errors",
			},
			[]string{"service", "error_type"},
		),
	}

	return metrics
}

// Record a request (by service and status code)
func (m *Metrics) RecordRequest(service string, statusCode int, durationSeconds float64) {

	// Convert status code to string (bucket of response response type)
	status := statusCodeToString(statusCode)

	m.RequestsTotal.WithLabelValues(service, status).Inc()

	m.RequestDuration.WithLabelValues(service).Observe(durationSeconds)
}

// Record an error (by service and type)
func (m *Metrics) RecordError(service string, errorType string) {

	m.ErrorsTotal.WithLabelValues(service, errorType).Inc()
}

// Increment the number of requests in flight
func (m *Metrics) IncInFlight() {
	m.RequestsInFlight.Inc()
}

// Decrement the number of requests in flight
func (m *Metrics) DecInFlight() {
	m.RequestsInFlight.Dec()
}

// Helper function to convert status code to string
func statusCodeToString(statusCode int) string {

	if statusCode >= 200 && statusCode < 300 {
		return "2xx"
	} else if statusCode >= 300 && statusCode < 400 {
		return "3xx"
	} else if statusCode >= 400 && statusCode < 500 {
		return "4xx"
	} else if statusCode >= 500 && statusCode < 600 {
		return "5xx"
	}
	return "unknown"
}