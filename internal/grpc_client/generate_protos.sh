#!/bin/bash

set -e

echo "Generating Golang protos..."

echo "Generating message proto..."
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    -I../../protos/ \
    ../../protos/prediction.proto

echo "Generating Golang protos... Done!"
