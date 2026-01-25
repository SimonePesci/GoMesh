# GoMesh

GoMesh is a lightweight HTTP reverse proxy written in Go. It sits in front of backend services, forwards traffic, and is the foundation for a small service mesh (data plane + control plane).

## Features

- HTTP reverse proxy with YAML configuration
- Graceful shutdown on SIGINT / SIGTERM
- Structured logging (Zap)
- Prometheus metrics on `/metrics`
- Panic recovery middleware (keeps the proxy running on handler panics)

## Requirements

- Go 1.24+

## Project layout

```
GoMesh/
├── cmd/
│   ├── proxy/
│   └── backend/
├── pkg/
│   ├── logging/
│   └── proxy/
├── config/proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

```bash
go mod download

go run cmd/backend/main.go
go run cmd/proxy/main.go

curl http://localhost:8000/api/users
curl http://localhost:8000/metrics

# intentional panic endpoint (proxy returns 500 and keeps running)
curl http://localhost:8000/panic
```

## Middleware stack

Requests pass through recovery, metrics, and logging middleware before the reverse proxy handler. Recovery is outermost so panics from any inner layer are caught.

## Configuration

Edit `config/proxy.yaml` for listen port, backend address, and timeouts.

## Build

```bash
make build
```

## License

See [LICENSE](LICENSE) when present.
