# GoMesh

GoMesh is a lightweight HTTP reverse proxy written in Go. It sits in front of backend services, forwards traffic, and is the foundation for a small service mesh (data plane + control plane).

## Features

- HTTP reverse proxy with YAML configuration
- Graceful shutdown on SIGINT / SIGTERM
- Structured logging (Zap)
- Prometheus metrics on `/metrics`
- Panic recovery middleware
- Distributed request tracing via `X-Trace-ID`

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
│   ├── tracing/    # Trace ID generation and helpers
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

curl -v http://localhost:8000/api/users
# Response includes X-Trace-ID

curl -H "X-Trace-ID: my-custom-trace-123" http://localhost:8000/api/users
curl http://localhost:8000/metrics
curl http://localhost:8000/panic
```

## Observability

- **Logs**: structured fields including `trace_id`, method, path, status, latency
- **Metrics**: `gomesh_requests_total`, `gomesh_request_duration_seconds`, `gomesh_requests_in_flight`, `gomesh_errors_total`
- **Tracing**: generate or accept `X-Trace-ID`, propagate to backends, return on the response

## Middleware stack

```
Recovery → Tracing → Metrics → Logging → ReverseProxy
```

## Configuration

Edit `config/proxy.yaml` for listen port, backend address, and timeouts.

## Build

```bash
make build
```

## License

MIT (see [LICENSE](LICENSE)).
