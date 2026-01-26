package tracing

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// Header name for the trace ID
const TraceIDHeader = "X-Trace-ID"

// Generate a new trace ID
// If the generation fails, it will return a fallback timestamp
func GenerateTraceID() string {
	
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {

		return fmt.Sprintf("fallback-%d", generateFallbackTimestamp())
	} 

	return hex.EncodeToString(bytes)
}


func generateFallbackTimestamp() int64 {
	return time.Now().UnixNano()
}

// Get the trace ID from the request header
func GetTraceID(r *http.Request) string {

	traceID := r.Header.Get(TraceIDHeader)
	if traceID == "" {
		return "unknown"
	}

	return traceID
}

// Set the trace ID in the request header
func SetTraceID(r *http.Request, traceID string) {
	r.Header.Set(TraceIDHeader, traceID)
}

// Set the trace ID in the response header
func SetTraceIDResponse(w http.ResponseWriter, traceID string) {
	w.Header().Set(TraceIDHeader, traceID)
}

