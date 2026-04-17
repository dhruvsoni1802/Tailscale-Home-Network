#!/bin/bash
set -e

echo "building React UI..."
cd ui && npm run build && cd ..

echo "building Go binaries..."
go build -o bin/server ./cmd/server
go build -o bin/client ./cmd/client

echo "done. binaries in bin/"