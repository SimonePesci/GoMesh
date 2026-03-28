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
	case sig := <-signChan:
		logger.Info("Received signal",
			zap.String("signal", sig.String()),
		)

		if err := server.Shutdown(10 * time.Second); err != nil {
			logger.Warn("Failed to shutdown server gracefully",
				zap.Error(err),
			)
		}
	}

	logger.Info("Proxy Terminated Successfully!")
}
