# GoMesh

GoMesh is a lightweight service mesh written in Go. It provides an HTTP data-plane proxy and a gRPC control plane that tracks proxies and streams routing configuration.

## Features

- HTTP reverse proxy with YAML bootstrap config
- Structured logging (Zap)
- Prometheus metrics (`/metrics`)
- Panic recovery middleware
- Request tracing (`X-Trace-ID`)
- gRPC control plane: proxy registration and config streaming
- Proxy identity and control-plane address in bootstrap config

## Requirements

- Go 1.24+

## Architecture

```
Client  →  Proxy (:8000)  →  Backend (:3000)
                ↑
         gRPC (register + stream)
                │
         Control plane (:9090)
```

## Project layout

```
GoMesh/
├── cmd/
│   ├── proxy/
│   ├── controller/
│   └── backend/
├── pkg/
│   ├── controlplane/
│   ├── logging/
│   ├── tracing/
│   └── proxy/          # includes control-plane client
├── api/proto/
├── scripts/generate-proto.sh
├── config/proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

```bash
go mod download

go run cmd/backend/main.go
go run cmd/controller/main.go
go run cmd/proxy/main.go

curl -v http://localhost:8000/api/users
curl http://localhost:8000/metrics
```

## Configuration

`config/proxy.yaml`:

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

Startup validates proxy identity, advertise address, and control-plane address.

## Control plane

- gRPC on port `9090` (`-port`, `-production`)
- `RegisterProxy` / `StreamConfig`
- Versioned route store with default backend `localhost:3000`

## Protocol Buffers

```bash
bash scripts/generate-proto.sh
```

## License

MIT. See [LICENSE](LICENSE).
