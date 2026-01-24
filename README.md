# GoMesh

GoMesh is a lightweight HTTP reverse proxy written in Go. It sits in front of backend services, forwards traffic, and is the foundation for a small service mesh (data plane + control plane).

## Features

- HTTP reverse proxy (`net/http`, `httputil.ReverseProxy`)
- YAML configuration
- Graceful shutdown on SIGINT / SIGTERM

## Requirements

- Go 1.24+

## Project layout

```
GoMesh/
├── cmd/
│   ├── proxy/      # Data plane proxy
│   └── backend/    # Test backend service
├── pkg/proxy/      # Proxy config, handler, server
├── config/
│   └── proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

### 1. Install dependencies

```bash
go mod download
```

### 2. Start the test backend

```bash
go run cmd/backend/main.go
```

Listens on `:3000`.

### 3. Start the proxy

```bash
go run cmd/proxy/main.go
```

Listens on `:8000` and forwards to the backend from `config/proxy.yaml`.

### 4. Send a request

```bash
curl http://localhost:8000/api/users
```

## Configuration

Edit `config/proxy.yaml`:

- Proxy listen port (default: `8000`)
- Backend host and port (default: `localhost:3000`)
- Read / write / idle timeouts

```bash
go run cmd/proxy/main.go -config /path/to/config.yaml
```

## Build

```bash
make build
# binaries in ./bin/
```

## License

See [LICENSE](LICENSE) when present.
