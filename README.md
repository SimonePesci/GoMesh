# GoMesh - Service Mesh Implementation

A high-performance, distributed Service Mesh with a gRPC Control Plane built from scratch in Go.

## Current Status: Phase 2 Part 4 - Distributed Tracing ✅

### What's Complete:

- ✅ **Phase 1**: Basic reverse proxy with graceful shutdown
- ✅ **Phase 2 Part 1**: Structured JSON logging with Zap and logging middleware
- ✅ **Phase 2 Part 2**: Prometheus metrics with `/metrics` endpoint
- ✅ **Phase 2 Part 3**: Recovery middleware and middleware chaining
- ✅ **Phase 2 Part 4**: Distributed tracing with unique trace IDs

### Next Up:

- 📍 **Phase 3**: gRPC Control Plane for dynamic configuration

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
│   ├── tracing/            # Distributed tracing (Phase 2 Part 4)
│   │   └── tracer.go       # Trace ID generation and propagation
│   └── proxy/              # Proxy package
│       ├── config.go       # Configuration loader
│       ├── handler.go      # Reverse proxy logic
│       ├── middleware.go   # All middleware (logging, metrics, tracing, recovery)
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

### Step 6: Test Recovery Middleware!

```bash
# Trigger a panic in the backend (proxy should handle it gracefully)
curl http://localhost:8000/panic

# You'll see:
# - Client receives: "Internal Server Error" (HTTP 500)
# - Proxy logs the panic with full stack trace
# - Proxy KEEPS RUNNING (doesn't crash!) ✅
```

### Step 7: Test Distributed Tracing!

```bash
# Make a request and check the trace ID in the response headers
curl -v http://localhost:8000/api/users

# You'll see in the response headers:
# X-Trace-ID: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6

# The response body also includes the trace ID:
# {
#   "message": "Hello from the backend service!",
#   "trace_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
#   ...
# }

# You can also provide your own trace ID:
curl -H "X-Trace-ID: my-custom-trace-123" http://localhost:8000/api/users

# The backend will echo it back in logs and response
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

- **trace_id**: Unique 128-bit identifier for tracking requests across services
- **method**: HTTP method (GET, POST, etc.)
- **path**: Request path
- **status**: HTTP status code (200, 404, 500, etc.)
- **latency_ms**: Time taken to process the request
- **timestamp**: ISO8601 formatted timestamp
- **remote_addr**: Client IP address

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

### Phase 2 Part 3: Advanced Middleware ✅

- ✅ **Panic Recovery** - Using `defer` + `recover()` to catch runtime panics
- ✅ **Stack Traces** - `runtime/debug.Stack()` for debugging panics
- ✅ **Resilient Design** - Isolated request failures don't crash the entire proxy
- ✅ **Middleware Chaining** - Variadic function pattern for composable middleware
- ✅ **Order-Aware Composition** - Reverse loop to apply middleware in correct order
- ✅ **Graceful Degradation** - Return 500 on panic, log details, keep serving

### Phase 2 Part 4: Distributed Tracing ✅

- ✅ **Trace ID Generation** - Using `crypto/rand` for cryptographically secure 128-bit IDs
- ✅ **Header Propagation** - `X-Trace-ID` header forwarded to backend services
- ✅ **Request Correlation** - Track requests across multiple services with unique IDs
- ✅ **Response Headers** - Return trace ID to clients for debugging
- ✅ **Logging Integration** - All logs include trace_id field for request correlation
- ✅ **Fallback Mechanism** - Timestamp-based fallback if crypto/rand fails

## Next Steps: Phase 3 - gRPC Control Plane

### Phase 3 Part 1: gRPC API Definition 📍 NEXT

- Define Protocol Buffer schemas for configuration
- Generate Go code from .proto files
- Design control plane API (GetConfig, UpdateConfig, etc.)

### Phase 3 Part 2: Control Plane Server

- Implement gRPC server for control plane
- Dynamic configuration management
- Hot reload without proxy restart

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

Understanding the request flow through our middleware chain:

```
Request from Client (with optional X-Trace-ID header)
        ↓
┌───────────────────────────┐
│  RecoveryMiddleware       │  ← OUTERMOST: Catches ALL panics (defer/recover)
│  (panic safety)           │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  TracingMiddleware        │  ← Generates/extracts trace ID
│  (distributed tracing)    │     Sets X-Trace-ID header in request & response
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  MetricsMiddleware        │  ← Increments in-flight counter, starts timer
│  (observability)          │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  LoggingMiddleware        │  ← Logs "request started" with trace_id
│  (structured logging)     │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  ProxyHandler             │  ← Forwards to backend with X-Trace-ID header
│  (httputil.ReverseProxy)  │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  Backend Service :3000    │  ← Receives trace ID, echoes in logs & response
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  LoggingMiddleware        │  ← Logs "request completed" with trace_id
│  (calculates latency)     │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  MetricsMiddleware        │  ← Records metrics (counter, histogram, gauge)
│  (records metrics)        │     Decrements in-flight counter
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  TracingMiddleware        │  ← X-Trace-ID already set in response header
│  (trace ID in response)   │
└───────────────────────────┘
        ↓
┌───────────────────────────┐
│  RecoveryMiddleware       │  ← If panic occurred, catches it here
│  (returns 500 if panic)   │     Logs with stack trace & trace_id, returns 500
└───────────────────────────┘
        ↓
Response to Client (with X-Trace-ID header)
```

### Middleware Order Matters!

The order of middleware is critical:
1. **Recovery** (outermost) - Must catch panics from ALL inner middleware
2. **Tracing** - Generate/extract trace ID early so all inner middleware can use it
3. **Metrics** - Track all requests (even if they panic)
4. **Logging** - Log all requests with trace ID (even if they panic)
5. **Proxy Handler** - The actual reverse proxy logic

### The Chain Function

Using the Chain helper makes middleware composition clean and readable:

```go
// Before: Manual nesting (hard to read)
handler = RecoveryMiddleware(logger,
    TracingMiddleware(
        MetricsMiddleware(metrics,
            LoggingMiddleware(logger, handler))))

// After: Chain function (clean & clear)
handler = Chain(handler,
    RecoveryMiddleware(logger),    // First = outermost
    TracingMiddleware(),           // Second
    MetricsMiddleware(metrics),    // Third
    LoggingMiddleware(logger),     // Last = innermost
)
```

**How Chain works:**

```go
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    // Apply in reverse so first middleware becomes outermost
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}
```

This produces: `Recovery(Tracing(Metrics(Logging(handler))))`

### Before vs After Logging

**Phase 1 (Basic Logging):**

```
[INFO] Forwarding: GET /api/users → localhost:3000
```

**Phase 2 Part 1-4 (Structured Logging + Tracing):**

```json
{"level":"info","timestamp":"2026-01-26T10:30:05.200Z","msg":"Request starter","trace_id":"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6","method":"GET","path":"/api/users","remote_addr":"127.0.0.1:53242"}
{"level":"info","timestamp":"2026-01-26T10:30:05.245Z","msg":"request completed","method":"GET","path":"/api/users","status":200,"latency_ms":"45ms","trace_id":"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"}
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

### Testing Recovery Middleware

The backend has a `/panic` endpoint that intentionally triggers a panic to test recovery:

```bash
# Trigger a panic
curl http://localhost:8000/panic
```

**What happens:**

1. **Client receives:**
   ```
   Internal Server Error
   ```
   HTTP status: 500

2. **Proxy logs show:**
   ```json
   {
     "level": "error",
     "msg": "panic recovered",
     "error": "intentional panic for testing recovery middleware!",
     "path": "/panic",
     "method": "GET",
     "stack": "goroutine 123 [running]:\nruntime/debug.Stack()..."
   }
   ```

3. **Proxy keeps running!** ✅
   - Other requests still work
   - No downtime
   - Panic isolated to single request

**Without recovery middleware, the entire proxy would crash.** This is the difference between:
- ❌ One bad request → entire proxy down → manual restart
- ✅ One bad request → 500 error → logged → proxy still running

### Distributed Tracing with Trace IDs

Every request gets a unique trace ID that follows it through your entire system:

```bash
# Make a request
curl -v http://localhost:8000/api/users
```

**1. Proxy receives request:**
```json
{"level":"info","msg":"Request starter","trace_id":"a1b2c3d4...","method":"GET","path":"/api/users"}
```

**2. Backend receives trace ID:**
```
[BACKEND] Received: GET /api/users | Trace-ID: a1b2c3d4...
```

**3. Backend returns response with trace ID:**
```json
{
  "message": "Hello from the backend service!",
  "trace_id": "a1b2c3d4...",
  ...
}
```

**4. Client receives trace ID in header:**
```
< X-Trace-ID: a1b2c3d4...
```

**Why This Matters:**

In a microservices architecture with multiple services:
```
Client → Proxy → Service A → Service B → Service C
```

**Without trace IDs:**
- ❌ Logs scattered across services
- ❌ Can't correlate which logs belong to same request
- ❌ Hard to debug issues spanning multiple services
- ❌ No visibility into request flow

**With trace IDs:**
- ✅ One ID tracks request through entire system
- ✅ Search logs for `trace_id: a1b2c3d4...` across all services
- ✅ See complete request journey
- ✅ Identify exactly where failures occurred
- ✅ Measure end-to-end latency

**Example debugging scenario:**
```bash
# User reports error, provides trace ID from response
# Search all logs for this trace ID:

# Proxy logs:
{"trace_id":"abc123","msg":"Request starter","path":"/api/order"}
{"trace_id":"abc123","msg":"request completed","status":500}

# Service A logs:
{"trace_id":"abc123","msg":"processing order","order_id":42}
{"trace_id":"abc123","msg":"calling payment service"}

# Service B logs:
{"trace_id":"abc123","error":"payment declined"}  ← Found the issue!
```

This is the foundation of observability in distributed systems!

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

### ✅ Phase 2 Part 3: Advanced Middleware (Complete)

- Middleware chaining helper with Chain() function
- Recovery middleware with panic handling
- Stack trace logging for debugging
- Resilient proxy that doesn't crash on panics

### ✅ Phase 2 Part 4: Distributed Tracing (Complete)

- Cryptographically secure trace ID generation (128-bit)
- X-Trace-ID header propagation to backend services
- Trace ID included in all logs for request correlation
- Response headers include trace ID for client debugging

### Phase 3: Control Plane with gRPC

- gRPC API definition
- Control plane server
- Dynamic configuration
- Hot reload

### Phase 4: Service Discovery & Load Balancing

- Service registry
- Round-robin load balancing
- Health checking

### Phase 5: Production Features

- mTLS encryption
- Circuit breaker
- Rate limiting
- JWT validation
