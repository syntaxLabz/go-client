package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "httpclient_requests_total",
			Help: "Total number of HTTP requests made",
		},
		[]string{"method", "status_code"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "httpclient_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status_code"},
	)
)

type metricsMiddleware struct {
	startTime time.Time
	method    string
}

// NewMetrics creates a new metrics middleware
func NewMetrics() Middleware {
	return &metricsMiddleware{}
}

func (m *metricsMiddleware) Before(req *http.Request) error {
	m.startTime = time.Now()
	m.method = req.Method
	return nil
}

func (m *metricsMiddleware) After(resp *http.Response) {
	duration := time.Since(m.startTime).Seconds()
	statusCode := strconv.Itoa(resp.StatusCode)

	requestsTotal.WithLabelValues(m.method, statusCode).Inc()
	requestDuration.WithLabelValues(m.method, statusCode).Observe(duration)
}