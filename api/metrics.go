package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	incomingEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "dasher",
			Subsystem: "http",
			Name:      "request_count",
			Help:      "Counter of requests received into the system.",
		}, []string{"service"})

	failedEvents = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "dasher",
			Subsystem: "http",
			Name:      "failed_total",
			Help:      "Counter of handle failures of requests (non-watches), by method (GET/PUT etc.) and code (400, 500 etc.).",
		}, []string{"method", "code"})

	successfulEventsHandlingTime = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "dasher",
			Subsystem: "http",
			Name:      "request_latency_microseconds",
			Help:      "Bucketed histogram of processing time (s) of successfully handled requests (non-watches), by service.",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 13),
		}, []string{"service"})
)

func init() {
	prometheus.MustRegister(incomingEvents)
	prometheus.MustRegister(failedEvents)
	prometheus.MustRegister(successfulEventsHandlingTime)
}

func requestCount(service string) {
	incomingEvents.WithLabelValues(service).Inc()
}

func requestLatency(service string, startTime time.Time) {
	successfulEventsHandlingTime.WithLabelValues(service).Observe(time.Since(startTime).Seconds())
}

func reportRequestFailed(request http.Request, err error) {
	method := request.Method
	failedEvents.WithLabelValues(method, strconv.Itoa(codeFromError(err))).Inc()
}

func codeFromError(err error) int {
	return http.StatusInternalServerError
}
