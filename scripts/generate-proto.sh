#!/bin/bash

# This script generates Go code from the Protocol Buffer definitions

echo "Generating Go code from proto files..."

protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    api/proto/mesh.proto

if [ $? -eq 0 ]; then
    echo "✅ Successfully generated:"
    echo "   - api/proto/mesh.pb.go        (message types)"
    echo "   - api/proto/mesh_grpc.pb.go   (gRPC service)"
else
    echo "❌ Generation failed!"
    echo ""
    echo "Make sure you have installed:"
    echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi