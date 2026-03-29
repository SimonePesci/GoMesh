# GoMesh

GoMesh is a lightweight service mesh written in Go. It provides an HTTP data-plane proxy and a gRPC control plane that registers proxies and streams versioned routing configuration.

## Features

- HTTP reverse proxy with YAML bootstrap config
- Structured logging (Zap)
- Prometheus metrics (`/metrics`)
- Panic recovery middleware
- Request tracing (`X-Trace-ID`)
- gRPC control plane (`RegisterProxy`, `StreamConfig`)
- Mandatory proxy registration and first config update before serving traffic
- In-memory storage of the latest streamed config on the proxy

## Requirements

- Go 1.24+

## Architecture

```
Client  →  Proxy (:8000)  →  Backend (:3000)
                ↑
         gRPC register + config stream
                │
         Control plane (:9090)
```

1. Proxy connects to the control plane and calls `RegisterProxy`
2. Proxy opens `StreamConfig` and waits for the initial `ConfigUpdate`
3. Proxy stores the config in memory and starts the HTTP server
4. If the config stream fails after startup, the proxy stops

Request forwarding still uses the static backend from bootstrap YAML until dynamic route application is implemented.

## Project layout

```
GoMesh/
├── cmd/
│   ├── proxy/              # Data plane
│   ├── controller/         # Control plane
│   └── backend/            # Test backend
├── pkg/
│   ├── controlplane/       # gRPC server, config store
│   ├── logging/            # Zap wrapper
│   ├── tracing/            # Trace IDs
│   └── proxy/              # HTTP server, middleware, CP client
├── api/proto/              # mesh.proto + generated stubs
├── scripts/generate-proto.sh
├── config/proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

```bash
go mod download

# 1. Backend
go run cmd/backend/main.go

# 2. Control plane
go run cmd/controller/main.go

# 3. Proxy (requires control plane)
go run cmd/proxy/main.go

# 4. Traffic
curl -v http://localhost:8000/api/users
curl http://localhost:8000/metrics
curl http://localhost:8000/panic
```

### Flags

**Controller**

| Flag | Default | Description |
|---|---|---|
| `-port` | `9090` | gRPC listen port |
| `-production` | `false` | JSON logging when true |

**Proxy**

| Flag | Default | Description |
|---|---|---|
| `-config` | `config/proxy.yaml` | Bootstrap config path |

## Configuration

```yaml
proxy:
  id: "proxy-1"
  version: "1.0.0"
  listen_port: 8000
  advertise_addr: "localhost:8000"
  backend:
    host: "localhost"
    port: 3000
  timeout:
    read_timeout: 30s
    write_timeout: 30s
    idle_timeout: 120s

control_plane:
  address: "localhost:9090"
```

## Observability

| Concern | Detail |
|---|---|
| Logging | Zap structured fields (`trace_id`, method, path, status, latency) |
| Metrics | `/metrics` — counters, histograms, in-flight gauge |
| Tracing | `X-Trace-ID` generated or accepted, propagated, returned |

### Middleware order

```
Recovery → Tracing → Metrics → Logging → ReverseProxy
```

## Control plane API

Defined in `api/proto/mesh.proto`:

| RPC | Type | Purpose |
|---|---|---|
| `RegisterProxy` | unary | Register a proxy instance |
| `StreamConfig` | server streaming | Push versioned route configs |

Regenerate stubs:

```bash
bash scripts/generate-proto.sh
```

## Roadmap

- Apply streamed config to request routing (replace static backend YAML)
- Service discovery and load balancing
- Production hardening (mTLS, circuit breaking, rate limits)

## Build

```bash
make build
```

## License

MIT. See [LICENSE](LICENSE).
