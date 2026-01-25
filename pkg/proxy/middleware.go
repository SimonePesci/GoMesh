package proxy

import (
	"net/http"
	"time"

	"github.com/SimonePesci/gomesh/pkg/logging"
	"go.uber.org/zap"
)

// this allows us to capture the status code of the response
// (the default ResponseWriter doesnt let you show the status code in the response)
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode: http.StatusOK,
		written: false,
	}
}

func (rw *responseWriter) WriteHeader(statusCode int) {

	// Only write the header if it hasn't been written yet
	// this is done to prevent multiple calls to WriteHeader()
	if !rw.written {
		rw.statusCode = statusCode
		rw.written = true
		rw.ResponseWriter.WriteHeader(statusCode)
	}

}

func (rw *responseWriter) Write(data []byte) (int, error) {
	// if the header hasn't been written yet, write it with the default status code
	// It must be written before calling Write()! (this is a Go requirement)
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(data)
}

func LoggingMiddleware(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record Start Time
		startTime := time.Now()

		wrappedWriter := newResponseWriter(w)

		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = "none"
		}


		logger.Info("Request starter",	
			zap.String("trace_id", traceID),
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.String("trace_id", traceID),
		)

		// call the next handler
		next.ServeHTTP(wrappedWriter, r)

		// Record Latency
		latency := time.Since(startTime)

		// Logs from the completed event
		logger.Info("request completed",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", wrappedWriter.statusCode),
			zap.Duration("latency_ms", latency),
			zap.String("trace_id", traceID),
		)
	})
}

