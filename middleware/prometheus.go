package AppMiddleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	helloRequestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "hello_request",
		Help:    "Histogram of the /hello request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
)

func WrapWithPrometheus(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		status := rw.status
		duration := time.Since(startTime).Seconds()
		helloRequestHistogram.WithLabelValues(
			r.Host+r.URL.String(),
			r.Method,
			fmt.Sprintf("%s %s", strconv.Itoa(status), http.StatusText(status)),
		).Observe(duration)
	}
	return http.HandlerFunc(fn)
}

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) Status() int {
	return rw.status
}
