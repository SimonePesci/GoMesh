# GoMesh

GoMesh is a lightweight service mesh written in Go. It provides an HTTP data-plane proxy and a gRPC control plane that tracks proxies and streams routing configuration.

## Features

- HTTP reverse proxy with YAML bootstrap config
- Structured logging (Zap)
- Prometheus metrics (`/metrics`)
- Panic recovery middleware
- Request tracing (`X-Trace-ID`)
- gRPC control plane with proxy registration and config streaming
- Mandatory control-plane connection and registration before the proxy serves traffic

## Requirements

- Go 1.24+

## Architecture

```
Client  →  Proxy (:8000)  →  Backend (:3000)
                ↑
         RegisterProxy + StreamConfig
                │
         Control plane (:9090)
```

The proxy fails fast at startup if the control plane is unreachable or registration fails.

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
│   └── proxy/
├── api/proto/
├── scripts/generate-proto.sh
├── config/proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

Start components in this order:

```bash
go mod download

go run cmd/backend/main.go
go run cmd/controller/main.go
go run cmd/proxy/main.go

curl -v http://localhost:8000/api/users
curl http://localhost:8000/metrics
curl http://localhost:8000/panic
```

### Controller flags

| Flag | Default | Description |
|---|---|---|
| `-port` | `9090` | gRPC listen port |
| `-production` | `false` | JSON logging when true |

### Proxy flags

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

Until dynamic routing is wired end-to-end, request forwarding still uses the static `backend` block. Registration with the control plane is already mandatory.

## Observability

- Structured logs with `trace_id`
- Prometheus at `/metrics`
- `X-Trace-ID` on requests and responses

## Middleware stack

```
Recovery → Tracing → Metrics → Logging → ReverseProxy
```

## Protocol Buffers

```bash
bash scripts/generate-proto.sh
```

## License

MIT. See [LICENSE](LICENSE).
