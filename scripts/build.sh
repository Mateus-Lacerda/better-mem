#!/bin/bash

base_path=$(dirname "$0")

mkdir -p ../target

echo "Building worker..."
go build -o $base_path/../target/worker.o $base_path/../cmd/worker/main.go

echo "Building api..."
go build -o $base_path/../target/api.o $base_path/../cmd/api/main.go
