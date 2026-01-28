package controlplane

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/SimonePesci/gomesh/api/proto"
	"go.uber.org/zap"
)

// Server is the control plane server
type Server struct {
	pb.UnimplementedMeshControlServer

	logger *zap.Logger
	configStore *ConfigStore

	mu sync.RWMutex
	proxies map[string]*ProxyConnection
}

// Represents a connection to a proxy: info and stream
type ProxyConnection struct {
	ProxyInfo *pb.ProxyInfo

	stream pb.MeshControl_StreamConfigServer
}


// NewServer creates a new control plane server
func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
		configStore: NewConfigStore(),
		proxies: make(map[string]*ProxyConnection),
	}
}

// Context is used to keep track of the context of the request (required by the grpc server)
func (s *Server) RegisterProxy(ctx context.Context, info *pb.ProxyInfo) (*pb.RegistrationResponse, error) {
	s.logger.Info("proxy registering",
		zap.String("proxy_id", info.ProxyId),
		zap.String("version", info.Version),
		zap.String("listen_addr", info.ListenAddr),
	)

	// Add the proxy to the map
	s.mu.Lock()
	s.proxies[info.ProxyId] = &ProxyConnection{
		ProxyInfo: info,
	}
	s.mu.Unlock()

	return &pb.RegistrationResponse{
		Success: true,
		Message: fmt.Sprintf("Proxy %s registered successfully!", info.ProxyId),
	}, nil
}

// StreamConfig is used to stream the config to the proxy
// Long-lived stream: server sends multiple messages
func (s *Server) StreamConfig(info *pb.ProxyInfo, stream pb.MeshControl_StreamConfigServer) error {
	s.logger.Info("proxy connecting for config stream",
		zap.String("proxy_id", info.ProxyId),
		zap.String("version", info.Version),
	)

	// Store the stream to send updates later
	s.mu.Lock()
	s.proxies[info.ProxyId] = &ProxyConnection{
		ProxyInfo: info,
		stream: stream,
	}
	s.mu.Unlock()

	// Remove proxy when connection closes
	defer func() {
		s.mu.Lock()
		delete(s.proxies, info.ProxyId)
		s.mu.Unlock()

		s.logger.Info("proxy disconnected",
			zap.String("proxy_id", info.ProxyId),
		)
	}()

	// Send the initial config
	config := s.configStore.GetConfig()
	s.logger.Info("sending initial config to proxy",
		zap.String("proxy_id", info.ProxyId),
		zap.Int64("version", config.Version),
		zap.Int("num_routes", len(config.Routes)),
	)

	if err := stream.Send(config); err != nil {
		s.logger.Error("failed to send initial config to proxy",
			zap.String("proxy_id", info.ProxyId),
			zap.Error(err),
		)
		return err
	}

	// We keep the connection alive
	// TODO: add the logic to handle config updates
	<- stream.Context().Done()

	return nil
}

// Broadcast update to all proxies
// Should be triggered by an admin when changing the configuration
func (s *Server) BroadcastConfigUpdate(config *pb.ConfigUpdate) {
	// No write lock, we just read the proxies map
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.logger.Info("broadcasting config update to proxies",
		zap.Int64("version", config.Version),
		zap.Int("proxy_count", len(s.proxies)),
	)

	// For each proxy, send the config update
	for proxyID, conn := range s.proxies {
		if conn.stream != nil {
			if err := conn.stream.Send(config); err != nil {
				s.logger.Error("failed to send config update to proxy",
					zap.String("proxy_id", proxyID),
					zap.Error(err),
				)
			} else {
				s.logger.Info("sent config update to proxy",
					zap.String("proxy_id", proxyID),
					zap.Int64("version", config.Version),
				)
			}
		}
	}

}

// GetConnectedProxies returns a list of all connected proxies
func (s *Server) GetConnectedProxies() []*pb.ProxyInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	proxies := make([]*pb.ProxyInfo, 0, len(s.proxies))
	for _, conn := range s.proxies {
		proxies = append(proxies, conn.ProxyInfo)
	}

	return proxies
}