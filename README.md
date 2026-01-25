# GoMesh

GoMesh is a lightweight HTTP reverse proxy written in Go. It sits in front of backend services, forwards traffic, and is the foundation for a small service mesh (data plane + control plane).

## Features

- HTTP reverse proxy with YAML configuration
- Graceful shutdown on SIGINT / SIGTERM
- Structured request logging with [Zap](https://github.com/uber-go/zap)

## Requirements

- Go 1.24+

## Project layout

```
GoMesh/
├── cmd/
│   ├── proxy/
│   └── backend/
├── pkg/
│   ├── logging/    # Zap logger wrapper
│   └── proxy/      # Config, handler, middleware, server
├── config/proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

```bash
go mod download

# terminal 1
go run cmd/backend/main.go

# terminal 2
go run cmd/proxy/main.go

# terminal 3
curl http://localhost:8000/api/users
```

Proxy logs include method, path, status, latency, and remote address as structured fields.

## Configuration

Edit `config/proxy.yaml` for listen port, backend address, and timeouts.

```bash
go run cmd/proxy/main.go -config /path/to/config.yaml
```

## Build

```bash
make build
```

## License

See [LICENSE](LICENSE) when present.
