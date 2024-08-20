#!/bin/bash

mkdir pkg/rpc/server/gen -p

poetry run python \
	-m grpc_tools.protoc \
	--proto_path=./proto \
	-Igen=./proto \
	--python_out=./pkg/rpc/server/gen \
	--pyi_out=./pkg/rpc/server/gen \
	--grpc_python_out=./pkg/rpc/server/gen \
	proto/*.proto

# Fix the imports in the generated Python files
for file in pkg/rpc/server/gen/*.py; do
	sed -i 's/^import SharedTypes_pb2/from . import SharedTypes_pb2/' $file
	sed -i 's/^import VDB_pb2/from . import VDB_pb2/' $file
	sed -i 's/^import Palace_pb2/from . import Palace_pb2/' $file
done
