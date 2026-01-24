package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// This is a simple backend service for testing the proxy

type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[BACKEND] Received: %s %s", r.Method, r.URL.Path)

	// Simulate some processing time
	time.Sleep(50 * time.Millisecond)

	// Create response
	resp := Response{
		Message:   "Hello from the backend service!",
		Timestamp: time.Now(),
		Path:      r.URL.Path,
		Method:    r.Method,
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

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/health", healthHandler)

	log.Println("[BACKEND] Starting test backend on :3000")
	log.Println("[BACKEND] Ready to receive requests from the proxy")
	
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalf("Backend failed: %v", err)
	}
}