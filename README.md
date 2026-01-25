# GoMesh - Service Mesh Implementation

A high-performance, distributed Service Mesh with a gRPC Control Plane built from scratch in Go.

## Current Status: Phase 2 Part 2 - Prometheus Metrics ✅

### What's Complete:

- ✅ **Phase 1**: Basic reverse proxy with graceful shutdown
- ✅ **Phase 2 Part 1**: Structured JSON logging with Zap and logging middleware
- ✅ **Phase 2 Part 2**: Prometheus metrics with `/metrics` endpoint

### Next Up:

- 📍 **Phase 2 Part 3**: Advanced middleware patterns (recovery, chaining)
- → **Phase 2 Part 4**: Distributed tracing with trace IDs

## Project Structure

```
GoMesh/
├── cmd/
│   ├── proxy/              # Main proxy binary
│   │   └── main.go
│   └── backend/            # Test backend service
│       └── main.go
├── pkg/
│   ├── logging/            # Structured logging (Phase 2 Part 1)
│   │   └── logging.go      # Zap logger wrapper
│   └── proxy/              # Proxy package
│       ├── config.go       # Configuration loader
│       ├── handler.go      # Reverse proxy logic
│       ├── middleware.go   # Logging & metrics middleware
│       ├── metrics.go      # Prometheus metrics (Phase 2 Part 2)
│       └── server.go       # HTTP server
├── config/
│   └── proxy.yaml          # Proxy configuration
├── go.mod
└── go.sum
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

You should see (structured JSON logs):

```json
{"level":"info","timestamp":"2026-01-25T...","msg":"Loading configuration from: config/proxy.yaml"}
{"level":"info","timestamp":"2026-01-25T...","msg":"GoMesh Proxy starting on port 8000"}
{"level":"info","timestamp":"2026-01-25T...","msg":"Forwarding all traffic to: http://localhost:3000"}
{"level":"info","timestamp":"2026-01-25T...","msg":"Press Ctrl+C to stop"}
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

### Step 5: Check Metrics!

```bash
# View Prometheus metrics
curl http://localhost:8000/metrics

# You'll see metrics like:
# gomesh_requests_total{service="backend",status="2xx"} 1
# gomesh_request_duration_seconds_bucket{service="backend",le="0.05"} 1
# gomesh_requests_in_flight 0
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

In the proxy terminal, you'll see structured JSON logs:

```json
{"level":"info","timestamp":"2026-01-25T10:30:05.200Z","msg":"Request starter","trace_id":"none","method":"GET","path":"/api/users","remote_addr":"127.0.0.1:53242"}
{"level":"info","timestamp":"2026-01-25T10:30:05.245Z","msg":"request completed","method":"GET","path":"/api/users","status":200,"latency_ms":"45ms","trace_id":"none"}
```

In the backend terminal, you'll see:

```
[BACKEND] Received: GET /api/users
```

### What the Logs Tell You:

- **trace_id**: Unique identifier for request tracing (will be used in Phase 2 Part 4)
- **method**: HTTP method (GET, POST, etc.)
- **path**: Request path
- **status**: HTTP status code (200, 404, 500, etc.)
- **latency_ms**: Time taken to process the request
- **timestamp**: ISO8601 formatted timestamp

## Configuration

Edit `config/proxy.yaml` to change:

- Proxy listen port (default: 8000)
- Backend host/port (default: localhost:3000)
- Timeouts

## What We've Learned

### Phase 1: Basic Reverse Proxy ✅

- ✅ **Go HTTP Server** - Using `net/http`
- ✅ **Reverse Proxy** - Using `httputil.ReverseProxy`
- ✅ **Configuration** - YAML parsing with `gopkg.in/yaml.v3`
- ✅ **Graceful Shutdown** - Signal handling with `os/signal`
- ✅ **Project Structure** - Standard Go project layout
- ✅ **Goroutines** - Running server in background
- ✅ **Channels** - Communication between goroutines

### Phase 2 Part 1: Structured Logging ✅

- ✅ **Zap Logger** - High-performance structured logging with `go.uber.org/zap`
- ✅ **Middleware Pattern** - Function wrapping: `func(http.Handler) http.Handler`
- ✅ **Interface Wrapping** - Custom `responseWriter` wraps `http.ResponseWriter`
- ✅ **Request/Response Tracking** - Capturing status codes and latency
- ✅ **Defer Pattern** - `defer logger.Sync()` ensures cleanup
- ✅ **Structured Fields** - JSON logs with typed fields (method, path, status, latency)

### Phase 2 Part 2: Prometheus Metrics ✅

- ✅ **Prometheus Client** - Using `github.com/prometheus/client_golang`
- ✅ **Metric Types** - Counters (requests_total), Histograms (duration), Gauges (in_flight)
- ✅ **Metric Labels** - Multi-dimensional metrics (service, status, error_type)
- ✅ **promauto Package** - Automatic registration with default registry
- ✅ **HTTP Multiplexer** - Using `http.NewServeMux()` for multiple endpoints
- ✅ **Middleware Chaining** - Metrics → Logging → Proxy handler stack
- ✅ **/metrics Endpoint** - Standard Prometheus scraping endpoint

## Next Steps: Phase 2 - Observability (Continued)

### Phase 2 Part 3: Advanced Middleware 📍 NEXT

- Middleware chaining helper
- Recovery middleware (panic handling)
- Prevent entire proxy from crashing

### Phase 2 Part 4: Distributed Tracing

- Generate unique trace IDs
- Inject trace IDs into request headers
- Track requests across multiple services

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

## How the Middleware Stack Works

Understanding the request flow through our middleware:

```
Request from Client
        ↓
┌───────────────────────────┐
│  MetricsMiddleware        │  ← Increments in-flight counter, starts timer
│  (wraps logging)          │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  LoggingMiddleware        │  ← Logs "request started"
│  (wraps handler)          │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  ProxyHandler             │  ← Forwards to backend (httputil.ReverseProxy)
│  (httputil.ReverseProxy)  │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  Backend Service :3000    │  ← Processes request, returns response
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  LoggingMiddleware        │  ← Logs "request completed" with status & latency
│  (calculates latency)     │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  MetricsMiddleware        │  ← Records metrics (counter, histogram, gauge)
│  (records metrics)        │     Decrements in-flight counter
└───────────────────────────┘
        ↓
Response to Client
```

### Before vs After Logging

**Phase 1 (Basic Logging):**

```
[INFO] Forwarding: GET /api/users → localhost:3000
```

**Phase 2 Part 1 (Structured Logging):**

```json
{"level":"info","timestamp":"2026-01-25T10:30:05.200Z","msg":"Request starter","trace_id":"none","method":"GET","path":"/api/users","remote_addr":"127.0.0.1:53242"}
{"level":"info","timestamp":"2026-01-25T10:30:05.245Z","msg":"request completed","method":"GET","path":"/api/users","status":200,"latency_ms":"45ms","trace_id":"none"}
```

The structured logs are:

- **Machine-readable** - Easy to parse and analyze
- **Searchable** - Query by any field (status=500, latency>100ms, etc.)
- **Standardized** - Consistent format across all services

### Prometheus Metrics

After sending some requests, check the `/metrics` endpoint:

```bash
curl http://localhost:8000/metrics
```

**Key Metrics Available:**

```prometheus
# Total number of requests (labeled by service and status bucket)
gomesh_requests_total{service="backend",status="2xx"} 5
gomesh_requests_total{service="backend",status="4xx"} 1
gomesh_requests_total{service="backend",status="5xx"} 0

# Request duration histogram (shows latency distribution)
gomesh_request_duration_seconds_bucket{service="backend",le="0.005"} 3
gomesh_request_duration_seconds_bucket{service="backend",le="0.01"} 5
gomesh_request_duration_seconds_bucket{service="backend",le="0.05"} 5
gomesh_request_duration_seconds_sum{service="backend"} 0.023
gomesh_request_duration_seconds_count{service="backend"} 5

# Number of requests currently being processed
gomesh_requests_in_flight 0

# Total errors (labeled by service and error type)
gomesh_errors_total{service="backend",error_type="timeout"} 0
```

**What These Metrics Tell You:**

- **requests_total**: Count requests by status code (2xx, 3xx, 4xx, 5xx)
- **request_duration_seconds**: Latency percentiles (p50, p95, p99) for SLA tracking
- **requests_in_flight**: Current load on the proxy
- **errors_total**: Error counts by type for alerting

---

## Complete Roadmap

### ✅ Phase 1: Basic Reverse Proxy (Complete)

- HTTP reverse proxy
- YAML configuration
- Graceful shutdown

### ✅ Phase 2 Part 1: Structured Logging (Complete)

- Zap logger integration
- Logging middleware
- Request/response tracking

### ✅ Phase 2 Part 2: Prometheus Metrics (Complete)

- Metrics package with Prometheus client
- `/metrics` endpoint for scraping
- Request counters, histograms, and gauges
- Metrics middleware for automatic tracking

### Phase 2 Part 3: Advanced Middleware

- Middleware chaining
- Recovery middleware

### Phase 2 Part 4: Distributed Tracing

- Trace ID generation
- Cross-service request tracking

### Phase 3: Control Plane with gRPC

- gRPC API definition
- Control plane server
- Dynamic configuration
- Hot reload

### Phase 4: Service Discovery & Load Balancing (~1

- Service registry
- Round-robin load balancing
- Health checking

### Phase 5: Production Features

- mTLS encryption
- Circuit breaker
- Rate limiting
- JWT validation
