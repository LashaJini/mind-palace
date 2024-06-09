.PHONY: all build deps dev-deps rpc clean-rpc rpc-py clean-rpc-py test-py test-go test db vdb godoc

BUILD_OUT_DIR=bin
BINARY_NAME=mind-palace
SOURCE_DIR=.

all: build

build: deps rpc
	@echo "Building the binary..."
	@go mod tidy
	@go build -o $(BUILD_OUT_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

deps:
	@go mod download
	@go mod verify

dev-deps:
	@go install --tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

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
	@go test -v ./cli/...

test: test-go test-py
	@echo "Done"

db:
	@bash scripts/postgres.sh $(ARGS)

vdb:
	@bash scripts/standalone_embed.sh $(ARGS)

godoc:
	@echo "Generating godoc..."
	@godoc -http=:6060 -play
