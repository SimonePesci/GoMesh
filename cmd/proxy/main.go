package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SimonePesci/gomesh/pkg/proxy"
)


func main() {

	configPath := flag.String("config", "config/proxy.yaml", "Path to config file")
	flag.Parse()

	log.Printf("[INFO] Loading configuration file from path: %s", *configPath)
	config, err := proxy.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("[FATAL] Failed to load config %v", err)
	}

	server, err := proxy.NewServer(config)
	if err != nil {
		log.Fatalf("[FATAL] Failed to create proxy server: %v", err)
	}

	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.Start()
	}()

	// Wait for shutdown: signal or error
	select {
	case err := <- serverErrors:
		if err != nil {
			log.Fatalf("[FATAL] Server error: %v", err)
		}
	case sig := <- signChan:
		log.Printf("[INFO] Received signal: %s, shutting down...", sig)

		if err := server.Shutdown(10 * time.Second); err != nil {
			log.Printf("[WARN] Failed to shutdown server gracefully: %v", err)
		}
	}


	log.Printf("[INFO] Proxy Terminated Successfully!")
}