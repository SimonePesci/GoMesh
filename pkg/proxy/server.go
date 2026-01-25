package proxy

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SimonePesci/gomesh/pkg/logging"
	"go.uber.org/zap"
)

type Server struct {
	config *Config
	handler *Handler
	httpServer *http.Server
	logger *logging.Logger
}

func NewServer(config *Config, logger *logging.Logger) (*Server, error) {

	// Create the handler
	handler, err := NewHandler(config, logger)
	if err != nil {
		return nil, fmt.Errorf("Failed to create handler for the server: %w", err)
	}

	// Wrap handler with logging middleware
	handlerWithMiddleware := LoggingMiddleware(logger, handler)

	// Create the http server for the proxy
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Proxy.ListenPort),
		Handler: handlerWithMiddleware,
		ReadTimeout: config.Proxy.Timeout.ReadTimeout,
		WriteTimeout: config.Proxy.Timeout.WriteTimeout,
		IdleTimeout: config.Proxy.Timeout.IdleTimeout,
	}

	return &Server{
		config: config,
		handler: handler,
		httpServer: httpServer,
		logger: logger,
	}, nil

}

// Starts the Server: will run till blocked
func (s *Server) Start() error {
	s.logger.Info("proxy server starting",
		zap.Int("port", s.config.Proxy.ListenPort),
		zap.String("backend_url", s.config.GetBackendURL()),
	)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error("failure in the server...stopping",
			zap.Error(err),
		)
		return fmt.Errorf("Failure in the server...stopping: %w", err)
	}

	return nil
}

// Handle Server closing gracefully
func (s *Server) Shutdown(timeout time.Duration) error {
	s.logger.Info("shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("Server shutdown failed: %w", err)
	}

	s.logger.Info("server stopped gracefully!")
	return nil
}