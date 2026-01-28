# GoMesh

GoMesh is a lightweight service mesh written in Go. It provides an HTTP data-plane proxy and a gRPC control plane that tracks proxies and streams routing configuration.

## Features

- HTTP reverse proxy with YAML bootstrap config
- Structured logging (Zap)
- Prometheus metrics (`/metrics`)
- Panic recovery middleware
- Request tracing (`X-Trace-ID`)
- gRPC control plane: proxy registration and config streaming
- Versioned in-memory route store on the control plane

## Requirements

- Go 1.24+

## Architecture

```
Client  →  Proxy (:8000)  →  Backend (:3000)
                ↑
         gRPC stream
                │
         Control plane (:9090)
```

## Project layout

```
GoMesh/
├── cmd/
│   ├── proxy/          # Data plane
│   ├── controller/     # Control plane
│   └── backend/        # Test backend
├── pkg/
│   ├── controlplane/   # gRPC server and config store
│   ├── logging/
│   ├── tracing/
│   └── proxy/
├── api/proto/          # mesh.proto and generated code
├── scripts/generate-proto.sh
├── config/proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

```bash
go mod download

# terminal 1 – test backend
go run cmd/backend/main.go

# terminal 2 – control plane
go run cmd/controller/main.go
# flags: -port (default 9090), -production

# terminal 3 – proxy
go run cmd/proxy/main.go

# terminal 4 – traffic
curl -v http://localhost:8000/api/users
curl http://localhost:8000/metrics
```

## Control plane

- Listens on gRPC port `9090` by default
- `RegisterProxy`: tracks connected proxies
- `StreamConfig`: sends versioned `ConfigUpdate` messages (default route to `localhost:3000`)
- Thread-safe config store with auto-incrementing versions
- Graceful shutdown on SIGINT / SIGTERM

## Protocol Buffers

```bash
bash scripts/generate-proto.sh
```

## Configuration

`config/proxy.yaml` controls proxy listen port, static backend, and timeouts.

## License

MIT. See [LICENSE](LICENSE).
