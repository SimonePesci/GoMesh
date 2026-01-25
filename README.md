# GoMesh

GoMesh is a lightweight HTTP reverse proxy written in Go. It sits in front of backend services, forwards traffic, and is the foundation for a small service mesh (data plane + control plane).

## Features

- HTTP reverse proxy with YAML configuration
- Graceful shutdown on SIGINT / SIGTERM
- Structured logging (Zap)
- Prometheus metrics on `/metrics`

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
│   └── proxy/      # Handler, middleware, metrics, server
├── config/proxy.yaml
├── Makefile
└── go.mod
```

## Quick start

```bash
go mod download

go run cmd/backend/main.go          # :3000
go run cmd/proxy/main.go            # :8000

curl http://localhost:8000/api/users
curl http://localhost:8000/metrics
```

### Metrics

| Metric | Type | Description |
|---|---|---|
| `gomesh_requests_total` | counter | Requests by service and status class |
| `gomesh_request_duration_seconds` | histogram | Request latency |
| `gomesh_requests_in_flight` | gauge | In-flight requests |
| `gomesh_errors_total` | counter | Errors by type |

## Configuration

Edit `config/proxy.yaml` for listen port, backend address, and timeouts.

## Build

```bash
make build
```

## License

See [LICENSE](LICENSE) when present.
