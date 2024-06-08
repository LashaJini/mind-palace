#!/bin/bash

mkdir pkg/rpc/server/gen -p

poetry run python \
	-m grpc_tools.protoc -Igen=./proto \
	--python_out=./pkg/rpc/server \
	--pyi_out=./pkg/rpc/server \
	--grpc_python_out=./pkg/rpc/server \
	proto/*.proto
