package proxy

import (
	"sync"

	pb "github.com/SimonePesci/gomesh/api/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// This manages the connection to the control plane
type ControlPlaneClient struct {
	logger *zap.Logger
	proxyID string
	version string
	listenAddr string

	controlPlaneAddr string
	conn *grpc.ClientConn
	client pb.MeshControlClient

	// current configuration
	mu sync.RWMutex
	config *pb.ConfigUpdate

	// callback for config updates
	onConfigUpdate func(config *pb.ConfigUpdate)
}