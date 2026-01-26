package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// This is a simple backend service for testing the proxy
// It simulates what "App B" would look like

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	TraceID   string    `json:"trace_id"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract trace ID from headers (forwarded by proxy)
	traceID := r.Header.Get("X-Trace-ID")
	if traceID == "" {
		traceID = "unknown"
	}
	
	log.Printf("[BACKEND] Received: %s %s | Trace-ID: %s", r.Method, r.URL.Path, traceID)

	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Create response
	resp := Response{
		Message:   "Hello from the backend service!",
		Timestamp: time.Now(),
		Path:      r.URL.Path,
		Method:    r.Method,
		TraceID:   traceID, 
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// panicHandler intentionally panics to test recovery middleware
func panicHandler(w http.ResponseWriter, r *http.Request) {
	traceID := r.Header.Get("X-Trace-ID")
	log.Printf("[BACKEND] Panic endpoint called - about to panic! Trace-ID: %s", traceID)
	panic("intentional panic for testing recovery middleware!")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/panic", panicHandler) // Test endpoint

	log.Println("[BACKEND] Starting test backend on :3000")
	log.Println("[BACKEND] Ready to receive requests from the proxy")
	log.Println("[BACKEND] Test panic recovery: curl http://localhost:8000/panic")
	log.Println("[BACKEND] Now echoing X-Trace-ID header in logs and responses")
	
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Backend failed: %v", err)
	}
}