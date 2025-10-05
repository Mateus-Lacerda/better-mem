#!/bin/bash

set -e

echo "Generating Python protos..."

echo "Generating message proto..."
uv run -m grpc_tools.protoc -I../protos --python_out=. --pyi_out=. --grpc_python_out=. ../protos/prediction.proto

echo "Generating Python protos... Done!"
