package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SimonePesci/gomesh/pkg/logging"
	"github.com/SimonePesci/gomesh/pkg/proxy"
	"go.uber.org/zap"
)

const controlPlaneStartupTimeout = 5 * time.Second

func main() {

	configPath := flag.String("config", "config/proxy.yaml", "Path to config file")
	production := flag.Bool("production", false, "Enable production mode (JSON logging)")
	flag.Parse()

	logger, err := logging.NewLogger(*production)
	if err != nil {
		panic("Failed to create logger: " + err.Error())
	}
	defer logger.Sync() // Flushes buffered log entries before exiting

	logger.Info("Loading configuration file from path",
		zap.String("path", *configPath),
	)
	config, err := proxy.LoadConfig(*configPath)
	if err != nil {
		logger.Fatal("Failed to load config",
			zap.Error(err),
		)
	}

	controlPlaneClient := proxy.NewControlPlaneClient(config, logger)

	connectCtx, cancelConnect := context.WithTimeout(context.Background(), controlPlaneStartupTimeout)
	if err := controlPlaneClient.Connect(connectCtx); err != nil {
		cancelConnect()
		logger.Fatal("Failed to connect to control plane",
			zap.String("address", config.ControlPlane.Address),
			zap.Error(err),
		)
	}
	cancelConnect()

	registerCtx, cancelRegister := context.WithTimeout(context.Background(), controlPlaneStartupTimeout)
	if _, err := controlPlaneClient.Register(registerCtx); err != nil {
		cancelRegister()

		if closeErr := controlPlaneClient.Close(); closeErr != nil {
			logger.Warn("Failed to close control plane client after registration error",
				zap.Error(closeErr),
			)
		}

		logger.Fatal("Failed to register proxy with control plane",
			zap.String("proxy_id", config.Proxy.ID),
			zap.Error(err),
		)
	}
	cancelRegister()

	streamCtx, cancelStream := context.WithCancel(context.Background())
	controlPlaneReady, controlPlaneErrors, err := controlPlaneClient.StartConfigStream(streamCtx)
	if err != nil {
		cancelStream()

		if closeErr := controlPlaneClient.Close(); closeErr != nil {
			logger.Warn("Failed to close control plane client after stream setup error",
				zap.Error(closeErr),
			)
		}

		logger.Fatal("Failed to start control plane config stream",
			zap.String("proxy_id", config.Proxy.ID),
			zap.Error(err),
		)
	}

	initialConfigCtx, cancelInitialConfig := context.WithTimeout(context.Background(), controlPlaneStartupTimeout)
	select {
	case <-controlPlaneReady:
		cancelInitialConfig()
	case err := <-controlPlaneErrors:
		cancelInitialConfig()
		cancelStream()

		if closeErr := controlPlaneClient.Close(); closeErr != nil {
			logger.Warn("Failed to close control plane client after initial config error",
				zap.Error(closeErr),
			)
		}

		logger.Fatal("Failed to receive initial config from control plane",
			zap.String("proxy_id", config.Proxy.ID),
			zap.Error(err),
		)
	case <-initialConfigCtx.Done():
		cancelInitialConfig()
		cancelStream()

		if closeErr := controlPlaneClient.Close(); closeErr != nil {
			logger.Warn("Failed to close control plane client after initial config timeout",
				zap.Error(closeErr),
			)
		}

		logger.Fatal("Timed out waiting for initial config from control plane",
			zap.String("proxy_id", config.Proxy.ID),
			zap.Duration("timeout", controlPlaneStartupTimeout),
		)
	}

	defer cancelStream()

	defer func() {
		if err := controlPlaneClient.Close(); err != nil {
			logger.Warn("Failed to close control plane client",
				zap.Error(err),
			)
		}
	}()

	server, err := proxy.NewServer(config, logger)
	if err != nil {
		logger.Fatal("Failed to create proxy server",
			zap.Error(err),
		)
	}

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.Start()
	}()

	// Wait for shutdown: signal or error
	select {
	case err := <-serverErrors:
		if err != nil {
			logger.Fatal("Server error",
				zap.Error(err),
			)
		}
	case err := <-controlPlaneErrors:
		cancelStream()
		logger.Fatal("Control plane config stream error",
			zap.String("proxy_id", config.Proxy.ID),
			zap.Error(err),
		)
	case sig := <-signChan:
		logger.Info("Received signal",
			zap.String("signal", sig.String()),
		)

		cancelStream()

		if err := server.Shutdown(10 * time.Second); err != nil {
			logger.Warn("Failed to shutdown server gracefully",
				zap.Error(err),
			)
		}
	}

	logger.Info("Proxy Terminated Successfully!")
}
