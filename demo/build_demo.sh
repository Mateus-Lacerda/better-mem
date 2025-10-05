#!/bin/bash

echo "Building BetterMem Demo..."
go build -o demo.o ./main.go

if [ $? -eq 0 ]; then
    echo "✓ Build successful!"
    echo "Run with: ./demo.o"
else
    echo "✗ Build failed!"
    exit 1
fi

