package proxy

import (
	"context"
	"fmt"
	"io"
	"sync"

	pb "github.com/SimonePesci/gomesh/api/proto"
	"github.com/SimonePesci/gomesh/pkg/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

// This manages the connection to the control plane.
type ControlPlaneClient struct {
	logger *logging.Logger

	proxyID       string
	version       string
	advertiseAddr string

	controlPlaneAddr string
	conn             *grpc.ClientConn
	client           pb.MeshControlClient

	// current configuration
	mu     sync.RWMutex
	config *pb.ConfigUpdate

	// callback for config updates
	onConfigUpdate func(config *pb.ConfigUpdate)
}

func NewControlPlaneClient(config *Config, logger *logging.Logger) *ControlPlaneClient {
	return &ControlPlaneClient{
		logger:           logger,
		proxyID:          config.Proxy.ID,
		version:          config.Proxy.Version,
		advertiseAddr:    config.Proxy.AdvertiseAddr,
		controlPlaneAddr: config.ControlPlane.Address,
	}
}

func (c *ControlPlaneClient) Connect(ctx context.Context) error {
	if c.conn != nil {
		return nil
	}

	c.logger.Info("connecting to control plane",
		zap.String("proxy_id", c.proxyID),
		zap.String("control_plane_addr", c.controlPlaneAddr),
	)

	conn, err := grpc.NewClient(
		c.controlPlaneAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to create control plane gRPC client: %w", err)
	}

	c.conn = conn
	c.client = pb.NewMeshControlClient(conn)
	c.conn.Connect()

	if err := c.waitUntilReady(ctx); err != nil {
		_ = c.Close()
		return err
	}

	c.logger.Info("control plane connection established",
		zap.String("proxy_id", c.proxyID),
		zap.String("control_plane_addr", c.controlPlaneAddr),
	)

	return nil
}

func (c *ControlPlaneClient) Register(ctx context.Context) (*pb.RegistrationResponse, error) {
	if c.client == nil {
		return nil, fmt.Errorf("control plane client is not connected")
	}

	info := c.proxyInfo()

	c.logger.Info("registering proxy with control plane",
		zap.String("proxy_id", info.ProxyId),
		zap.String("version", info.Version),
		zap.String("listen_addr", info.ListenAddr),
	)

	response, err := c.client.RegisterProxy(ctx, info)
	if err != nil {
		return nil, fmt.Errorf("failed to register proxy with control plane: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("control plane rejected proxy registration: %s", response.Message)
	}

	c.logger.Info("proxy registered with control plane",
		zap.String("proxy_id", info.ProxyId),
		zap.String("message", response.Message),
	)

	return response, nil
}

func (c *ControlPlaneClient) StartConfigStream(ctx context.Context) (<-chan struct{}, <-chan error, error) {
	if c.client == nil {
		return nil, nil, fmt.Errorf("control plane client is not connected")
	}

	c.logger.Info("subscribing to control plane config stream",
		zap.String("proxy_id", c.proxyID),
	)

	stream, err := c.client.StreamConfig(ctx, c.proxyInfo())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start control plane config stream: %w", err)
	}

	c.logger.Info("control plane config stream established",
		zap.String("proxy_id", c.proxyID),
	)

	readyCh := make(chan struct{})
	errCh := make(chan error, 1)

	go c.receiveConfigUpdates(ctx, stream, readyCh, errCh)

	return readyCh, errCh, nil
}

func (c *ControlPlaneClient) GetConfig() *pb.ConfigUpdate {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return cloneConfig(c.config)
}

func (c *ControlPlaneClient) Close() error {
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	c.conn = nil
	c.client = nil

	if err != nil {
		return fmt.Errorf("failed to close control plane connection: %w", err)
	}

	c.logger.Info("control plane client connection closed",
		zap.String("proxy_id", c.proxyID),
	)

	return nil
}

func (c *ControlPlaneClient) proxyInfo() *pb.ProxyInfo {
	return &pb.ProxyInfo{
		ProxyId:    c.proxyID,
		Version:    c.version,
		ListenAddr: c.advertiseAddr,
	}
}

func (c *ControlPlaneClient) receiveConfigUpdates(
	ctx context.Context,
	stream pb.MeshControl_StreamConfigClient,
	readyCh chan<- struct{},
	errCh chan<- error,
) {
	receivedInitialConfig := false

	for {
		update, err := stream.Recv()
		if err != nil {
			if ctx.Err() != nil {
				return
			}

			if err == io.EOF {
				errCh <- fmt.Errorf("control plane config stream closed unexpectedly")
				return
			}

			errCh <- fmt.Errorf("failed to receive config update from control plane: %w", err)
			return
		}

		config := c.storeConfig(update)

		c.logger.Info("received config update from control plane",
			zap.String("proxy_id", c.proxyID),
			zap.Int64("version", config.Version),
			zap.Int("num_routes", len(config.Routes)),
		)

		if !receivedInitialConfig {
			receivedInitialConfig = true
			close(readyCh)
		}

		if c.onConfigUpdate != nil {
			c.onConfigUpdate(cloneConfig(config))
		}
	}
}

func (c *ControlPlaneClient) waitUntilReady(ctx context.Context) error {
	for {
		state := c.conn.GetState()

		switch state {
		case connectivity.Ready:
			return nil
		case connectivity.Shutdown:
			return fmt.Errorf("control plane connection shut down before becoming ready")
		}

		if !c.conn.WaitForStateChange(ctx, state) {
			return fmt.Errorf("timed out waiting for control plane connection to become ready: %w", ctx.Err())
		}
	}
}

func (c *ControlPlaneClient) storeConfig(config *pb.ConfigUpdate) *pb.ConfigUpdate {
	clonedConfig := cloneConfig(config)

	c.mu.Lock()
	c.config = clonedConfig
	c.mu.Unlock()

	return cloneConfig(clonedConfig)
}

func cloneConfig(config *pb.ConfigUpdate) *pb.ConfigUpdate {
	if config == nil {
		return nil
	}

	return proto.Clone(config).(*pb.ConfigUpdate)
}
