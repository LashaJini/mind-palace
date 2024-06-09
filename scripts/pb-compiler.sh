#!/bin/bash

if ! which protoc >/dev/null; then
	echo "Please install protocol buffer compiler: https://grpc.io/docs/protoc-installation"
	exit 1
fi

if ! which protoc-gen-go >/dev/null; then
	echo -e "Please install protoc-gen-go:\n"
	echo "go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0"
	exit 1
fi

if ! which protoc-gen-go-grpc >/dev/null; then
	echo -e "Please install protoc-gen-go-grpc:\n"
	echo "go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0"
	exit 1
fi

mkdir pkg/rpc/client/gen -p

protoc \
	--go_out=./pkg/rpc/client/gen --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/rpc/client/gen --go-grpc_opt=paths=source_relative \
	proto/*.proto

bash ./scripts/pb-compiler-py.sh
