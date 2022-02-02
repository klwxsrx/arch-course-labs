package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

type HTTPHandler interface {
	MetricsHandler() http.Handler
	Middleware(excludedEndpoints []string) func(next http.Handler) http.Handler
}

type httpHandler struct {
	latencyHistogram *prometheus.HistogramVec
	requestCounter   *prometheus.CounterVec
}

func (h *httpHandler) MetricsHandler() http.Handler {
	return promhttp.Handler()
}

func (h *httpHandler) Middleware(excludedEndpoints []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if h.hasSpecifiedEndpoint(r.RequestURI, excludedEndpoints) {
				next.ServeHTTP(w, r)
				return
			}

			lrw := newLoggingResponseWriter(w)

			now := time.Now()
			next.ServeHTTP(lrw, r)
			duration := time.Since(now)

			h.addLatencyStat(r.Method, r.RequestURI, duration, lrw.code)
			h.incRequestsCount(r.Method, r.RequestURI, lrw.code)
		})
	}
}

func (h *httpHandler) addLatencyStat(method, endpoint string, duration time.Duration, httpStatus int) {
	h.latencyHistogram.WithLabelValues(method, endpoint, strconv.Itoa(httpStatus)).Observe(duration.Seconds())
}

func (h *httpHandler) incRequestsCount(method, endpoint string, httpStatus int) {
	h.requestCounter.WithLabelValues(method, endpoint, strconv.Itoa(httpStatus)).Inc()
}

func (h *httpHandler) hasSpecifiedEndpoint(endpoint string, expectedEndpoints []string) bool {
	for _, expected := range expectedEndpoints {
		if endpoint == expected {
			return true
		}
	}
	return false
}

func NewHTTPHandler() (HTTPHandler, error) {
	latencyHistogram := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "app_request_latency_seconds",
		Help: "Application Request Latency",
	}, []string{"method", "endpoint", "http_status"})
	err := prometheus.Register(latencyHistogram)
	if err != nil {
		return nil, err
	}

	requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "app_request_count",
		Help: "Application Request Count",
	}, []string{"method", "endpoint", "http_status"})
	err = prometheus.Register(requestCounter)
	if err != nil {
		return nil, err
	}

	return &httpHandler{
		latencyHistogram: latencyHistogram,
		requestCounter:   requestCounter,
	}, nil
}

type loggingResponseWriter struct {
	http.ResponseWriter
	code int
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}
