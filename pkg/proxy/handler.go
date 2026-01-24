package proxy

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// Proxy struct, a reverse proxy reference and a config reference
type Handler struct {
	config *Config
	reverseProxy *httputil.ReverseProxy
}

// Builds a new Handler
func NewHandler(config *Config) (*Handler, error) {

	// Parse the Backend URL from Config file
	rawBackendURL := config.GetBackendURL()
	backendURL, err := url.Parse(rawBackendURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse backend URL, is it written correctly?")
	}

	// Create a new reverse proxy from the builtin Go lib (it copies headers and streams)
	reverseProxy := httputil.NewSingleHostReverseProxy(backendURL)

	// Customize proxy to handle errors differently
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[ERROR] Proxy error: %v | URL: %s", err, r.URL.Path)
		http.Error(w, "Gateway Error", http.StatusBadGateway)
	}

	// Modify outgoing requests to backend
	originalDirector := reverseProxy.Director
	originalDirector = func(req *http.Request) {

		originalDirector(req)

		req.Header.Set("X-Forwarded-By", "GoMesh-Proxy")

		log.Printf("Forwarding %s %s -> %s", req.Method, req.URL, backendURL)
	}


	return &Handler{
		config: config,
		reverseProxy: reverseProxy,
	}, nil

}


// Serve through the reverse Proxy
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: To implement later
	// routing logic
	// load balancing
	// Circuit Breaking 
	// Rate limiting

	h.reverseProxy.ServeHTTP(w, r)
}