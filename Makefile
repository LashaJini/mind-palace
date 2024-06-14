.PHONY: all build deps deps-go deps-py dev-deps rpc clean-rpc rpc-py clean-rpc-py test-py test-go test cover db vdb graph godoc

BUILD_OUT_DIR=bin
BINARY_NAME=mind-palace
SOURCE_DIR=.

.EXPORT_ALL_VARIABLES:
PROJECT_ROOT=$(shell pwd)
MP_ENV ?= dev# prod,test,dev

all: build

build: deps rpc
	@echo "Building the binary..."
	@go mod tidy
	@go build -o $(BUILD_OUT_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

deps: deps-go deps-py

deps-go:
	@echo "Installing go dependencies..."
	@go mod download
	@go mod verify

deps-py:
	@echo "Installing python dependencies..."
	@poetry install

dev-deps:
	@go install github.com/kisielk/godepgraph@latest
	@go install golang.org/x/tools/cmd/cover@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	@go install --tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@poetry add pytest-cov

rpc:
	@echo "Compiling '.proto' files..."
	@bash ./scripts/pb-compiler.sh

clean-rpc:
	@echo "Removing compiled '.proto' files..."
	@rm -rf ./pkg/rpc/client/gen ./pkg/rpc/server/gen

rpc-py:
	@echo "Compiling '.proto' files..."
	@bash ./scripts/pb-compiler-py.sh

clean-rpc-py:
	@echo "Removing compiled '.proto' python files..."
	@rm -rf ./pkg/rpc/server/gen

test-py:
	@poetry run pytest

test-go:
	MP_ENV=test LOG_LEVEL=5 go test -v $(shell go list ./...) $(ARGS)

test: test-go test-py
	@echo "Done"

cover: dev-deps
	MP_ENV=test LOG_LEVEL=5 bash scripts/cover.sh $(ARGS)

db:
	@bash scripts/postgres.sh $(ARGS)

vdb:
	@bash scripts/standalone_embed.sh $(ARGS)

graph:
	@godepgraph -p \
		google,github.com/google,github.com/lib,github.com/joho,github.com/spf13 \
		-stoponerror=false \
		-s . | dot -Tpng -o godepgraph.png
	@eog godepgraph.png

godoc:
	@echo "Generating godoc..."
	@godoc -http=:6060 -play
