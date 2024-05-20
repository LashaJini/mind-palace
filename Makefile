.PHONY: all build deps dev-deps rpc clean-rpc db vdb godoc

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
	@echo "Removing '.proto' files..."
	@rm -rf ./pkg/rpc/client/gen ./pkg/rpc/server/gen

db:
	@bash scripts/postgres.sh $(ARGS)

vdb:
	@bash scripts/standalone_embed.sh $(ARGS)

godoc:
	@echo "Generating godoc..."
	@godoc -http=:6060 -play
