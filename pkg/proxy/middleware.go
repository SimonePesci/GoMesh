package proxy

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/SimonePesci/gomesh/pkg/logging"
	"github.com/SimonePesci/gomesh/pkg/tracing"
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

func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceID := tracing.GetTraceID(r)
		// If the trace ID is unknown, it means it's the first request
		// so we generate a new trace ID and set it in the request header
		if traceID == "unknown" {
			traceID = tracing.GenerateTraceID()
			tracing.SetTraceID(r, traceID)
		}

		// Set the trace ID in the response header
		// So the client can use it to trace the request
		tracing.SetTraceIDResponse(w, traceID)

		// Call the next handler
		next.ServeHTTP(w, r)

	})
}

func LoggingMiddleware(logger *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record Start Time
		startTime := time.Now()

		wrappedWriter := newResponseWriter(w)

		traceID := tracing.GetTraceID(r)


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

// Middleware to record metrics for the request
func MetricsMiddleware(metrics *Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Avoiding metrics in metrics: would cause infinite recursion!
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()

		metrics.IncInFlight()
		defer metrics.DecInFlight() // this will ensure decrementing the in flight counter even with a panic

		wrappedWriter := newResponseWriter(w)

		next.ServeHTTP(wrappedWriter, r)

		// In Seconds to be compatible with Prometheus (which uses seconds for the histogram)
		duration := time.Since(startTime).Seconds()

		// Record the request metrics
		// TODO: get the service name from the request header
		metrics.RecordRequest("backend", wrappedWriter.statusCode, duration)
	})
}

// Middleware to recover from panics and log the error
// This will prevent the entire proxy from crashing
func RecoveryMiddleware(logging *logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Defer the recovery
		// it will be executed after the next.ServeHTTP() call
		defer func() {
			// Log the panic with stack trace
			if err := recover(); err != nil {


				traceID := tracing.GetTraceID(r)

				logging.Error("panic recovered",
					zap.String("trace_id", traceID),
					zap.Any("error", err),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.String("stack", string(debug.Stack())),
				)

				// Return a 500 Internal Server Error to the client
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)

	})
}

// Middleware chainer
// This will apply middlewares in the order they appear in the list
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {

	// Apply middlewares in reverse order so the first one becomes outermost
	// This ensures: Chain(h, A, B, C) produces: A(B(C(h)))
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}

