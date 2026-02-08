#!/bin/bash

if ! command -v air &> /dev/null
then
    echo "air could not be found"
    echo "check out https://github.com/air-verse/air"
    exit
fi

if [ ! -d "bin" ]; then
    mkdir bin
fi

air --build.cmd "go build -tags=local -o bin/api cmd/api" --build.bin "./bin/api"
