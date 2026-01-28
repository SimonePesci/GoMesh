package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/SimonePesci/gomesh/pkg/controlplane"

	pb "github.com/SimonePesci/gomesh/api/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)


func main() {

	port := flag.Int("port", 9090, "Port the server will listen on for gRPC connections")
	production := flag.Bool("production", false, "Whether to run in production mode (JSON logging)")
	flag.Parse()

	var logger *zap.Logger
	var err error

	if *production {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}
	defer logger.Sync() // to flush logs before exiting

	logger.Info("Control Plane starting...",
		zap.Int("port", *port),
		zap.Bool("production", *production),
	)

	// Create the control plane server
	controlPlane := controlplane.NewServer(logger)

	// Create the gRPC server
	grpcServer := grpc.NewServer()

	// Now register both
	pb.RegisterMeshControlServer(grpcServer, controlPlane)

	logger.Info("gRPC server registered")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logger.Fatal("failed to start server",
			zap.Int("Port:", *port),
			zap.Error(err),
		)
	}

	logger.Info("control plane listening at",
		zap.String("address", listener.Addr().String()),
	)

	// shudown procedure
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		// We pick up errors from the grpc server
		serverErrors <- grpcServer.Serve(listener)
	}()

	select {

	// Handle errors coming from the grpc server
	case err := <- serverErrors:
		if err != nil {
			logger.Fatal("server error",
				zap.Error(err),
			)
		}

	// Handle signals (e.g. Ctrl+C)
	case sig := <- sigChan:
		logger.Info("received signal",
			zap.String("signal", sig.String()),
		)

		logger.Info("shutting down server gracefully...")
		grpcServer.GracefulStop()
		logger.Info("server terminated gracefully")
	}

	logger.Info("control plane terminated successfully")
}