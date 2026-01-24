.PHONY: help run-backend run-proxy test build clean

help:
	@echo "GoMesh - Available Commands:"
	@echo "  make run-backend    - Start the test backend service"
	@echo "  make run-proxy      - Start the proxy"
	@echo "  make test          - Test the proxy with curl"
	@echo "  make build         - Build both binaries"
	@echo "  make clean         - Remove built binaries"

run-backend:
	@echo "Starting backend service on :3000..."
	go run cmd/backend/main.go

run-proxy:
	@echo "Starting proxy on :8000..."
	go run cmd/proxy/main.go

test:
	@echo "Testing proxy..."
	@echo "\n1. Health check:"
	curl -s http://localhost:8000/health | jq || curl -s http://localhost:8000/health
	@echo "\n\n2. API request:"
	curl -s http://localhost:8000/api/users | jq || curl -s http://localhost:8000/api/users
	@echo "\n\n3. POST request:"
	curl -s -X POST http://localhost:8000/api/data -d '{"test":"data"}' | jq || curl -s -X POST http://localhost:8000/api/data -d '{"test":"data"}'

build:
	@echo "Building binaries..."
	@mkdir -p bin
	go build -o bin/proxy cmd/proxy/main.go
	go build -o bin/backend cmd/backend/main.go
	@echo "✓ Binaries created in ./bin/"

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	@echo "✓ Clean complete"