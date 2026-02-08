#!/bin/bash

base_path=$(dirname "$0")
bin_path=$base_path/../bin
mkdir -p $bin_path
cmd_path=$base_path/../cmd

echo "Building worker..."
go build -tags=local -o $bin_path/worker $cmd_path/worker

echo "Building api..."
go build -tags=local -o $bin_path/api $cmd_path/api
