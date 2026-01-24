package proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	config *Config
	handler *Handler
	httpServer *http.Server
}

func NewServer(config *Config) (*Server, error) {

	// Create the handler
	handler, err := NewHandler(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create handler for the server: %w", err)
	}

	// Create the http server for the proxy
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Proxy.ListenPort),
		Handler: handler,
		ReadTimeout: config.Proxy.Timeout.ReadTimeout,
		WriteTimeout: config.Proxy.Timeout.WriteTimeout,
		IdleTimeout: config.Proxy.Timeout.IdleTimeout,
	}

	return &Server{
		config: config,
		handler: handler,
		httpServer: httpServer,
	}, nil

}

// Starts the Server: will run till blocked
func (s *Server) Start() error {
	log.Printf("[INFO] GoMesh Proxy starting at port: %d", s.config.Proxy.ListenPort)
	log.Printf("[INFO] Forwarding all traffic to %s", s.config.GetBackendURL())
	log.Printf("[INFO] Use Ctrl+C to stop server")

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("Failure in the server...stopping: %w", err)
	}

	return nil
}

// Handle Server closing gracefully
func (s *Server) Shutdown(timeout time.Duration) error {
	log.Printf("[INFO] Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server shutdown failed: %w", err)
	}

	log.Printf("[INFO] Server stopped gracefully!")
	return nil
}