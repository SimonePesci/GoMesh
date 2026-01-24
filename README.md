# GoMesh - Phase 1: Basic Reverse Proxy

A high-performance, distributed Service Mesh with a gRPC Control Plane.

## Current Status: Phase 1 - Basic Proxy ✅

We've built a working reverse proxy that forwards HTTP traffic to a backend service.

## Project Structure

```
gomesh/
├── cmd/
│   ├── proxy/           # Main proxy binary
│   │   └── main.go
│   └── backend/         # Test backend service
│       └── main.go
├── pkg/
│   └── proxy/           # Proxy package
│       ├── config.go    # Configuration loader
│       ├── handler.go   # Reverse proxy logic
│       └── server.go    # HTTP server
├── config/
│   └── proxy.yaml       # Proxy configuration
└── go.mod
```

## How to Run

### Step 1: Install Dependencies

```bash
cd gomesh
go mod download
```

### Step 2: Start the Backend Service

In one terminal:

```bash
go run cmd/backend/main.go
```

You should see:

```
[BACKEND] Starting test backend on :3000
[BACKEND] Ready to receive requests from the proxy
```

### Step 3: Start the Proxy

In another terminal:

```bash
go run cmd/proxy/main.go
```

You should see:

```
[INFO] Loading configuration from: config/proxy.yaml
[INFO] GoMesh Proxy starting on port 8000
[INFO] Forwarding all traffic to: http://localhost:3000
[INFO] Press Ctrl+C to stop
```

### Step 4: Test It!

In a third terminal:

```bash
# Send a request to the proxy
curl http://localhost:8000/api/users

# You should see:
# {
#   "message": "Hello from the backend service!",
#   "timestamp": "2026-01-24T...",
#   "path": "/api/users",
#   "method": "GET"
# }
```

## What's Happening?

```
Your Request
    ↓
    ↓ (HTTP GET /api/users)
    ↓
┌─────────────────┐
│  GoMesh Proxy   │  ← Listening on :8000
│  (localhost)    │
└─────────────────┘
    ↓
    ↓ (Forwards to backend)
    ↓
┌─────────────────┐
│  Backend App    │  ← Listening on :3000
│  (localhost)    │
└─────────────────┘
    ↓
    ↓ (Returns JSON)
    ↓
Your Response
```

## Watch the Logs

In the proxy terminal, you'll see:

```
[INFO] Forwarding: GET /api/users → localhost:3000
```

In the backend terminal, you'll see:

```
[BACKEND] Received: GET /api/users
```

## Configuration

Edit `config/proxy.yaml` to change:

- Proxy listen port (default: 8000)
- Backend host/port (default: localhost:3000)
- Timeouts

## What We've Learned (Phase 1)

✅ **Go HTTP Server** - Using `net/http`
✅ **Reverse Proxy** - Using `httputil.ReverseProxy`
✅ **Configuration** - YAML parsing with `gopkg.in/yaml.v3`
✅ **Graceful Shutdown** - Signal handling with `os/signal`
✅ **Project Structure** - Standard Go project layout

## Next Steps: Phase 2 - Observability

In the next phase, we'll add:

- Structured logging (Zap)
- Prometheus metrics (`/metrics` endpoint)
- Request tracing (Trace-ID injection)
- Middleware pattern

## Troubleshooting

**Port already in use?**

```bash
# Find what's using port 8000
lsof -i :8000

# Or change the port in config/proxy.yaml
```

**Can't connect to backend?**

```bash
# Make sure backend is running
curl http://localhost:3000/health
```

## Commands Cheat Sheet

```bash
# Run proxy
go run cmd/proxy/main.go

# Run proxy with custom config
go run cmd/proxy/main.go -config /path/to/config.yaml

# Build the binary
go build -o bin/proxy cmd/proxy/main.go

# Run the binary
./bin/proxy

# Test the proxy
curl -v http://localhost:8000/test
curl -X POST http://localhost:8000/api/data -d '{"key":"value"}'
```

---

**You're ready for Phase 1!** 🚀

Try running it and experimenting with different requests. When you're comfortable, we'll move to Phase 2: Observability.
