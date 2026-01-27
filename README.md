# GoMesh

GoMesh is a lightweight service mesh written in Go. The data plane is an HTTP reverse proxy with logging, metrics, recovery, and tracing. The control plane API is defined with gRPC and Protocol Buffers for proxy registration and configuration streaming.

## Features

- HTTP reverse proxy with YAML configuration
- Structured logging (Zap)
- Prometheus metrics (`/metrics`)
- Panic recovery middleware
- Request tracing (`X-Trace-ID`)
- gRPC API contract for control plane ↔ data plane communication

## Requirements

- Go 1.24+
- Optional: `protoc` and Go plugins to regenerate stubs

## Project layout

```
GoMesh/
├── cmd/
│   ├── proxy/
│   └── backend/
├── pkg/
│   ├── logging/
│   ├── tracing/
│   └── proxy/
├── api/proto/
│   ├── mesh.proto
│   ├── mesh.pb.go
│   └── mesh_grpc.pb.go
├── scripts/generate-proto.sh
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
curl http://localhost:8000/metrics
```

## gRPC API

`api/proto/mesh.proto` defines `MeshControl`:

| RPC | Type | Purpose |
|---|---|---|
| `RegisterProxy` | unary | Proxy registers with the control plane |
| `StreamConfig` | server streaming | Control plane pushes config updates |

Regenerate code after editing the proto:

```bash
# protoc + protoc-gen-go + protoc-gen-go-grpc required
bash scripts/generate-proto.sh
```

Do not edit the generated `.pb.go` files by hand.

## Observability

- Structured logs with `trace_id`
- Prometheus endpoint at `/metrics`
- `X-Trace-ID` propagation

## Configuration

Edit `config/proxy.yaml` for listen port, backend address, and timeouts.

## License

MIT. See [LICENSE](LICENSE).
