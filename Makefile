.PHONY: all build godoc help

BUILD_OUT_DIR=bin
BINARY_NAME=mind-palace
SOURCE_DIR=.

all: build

build: deps
	@echo "Building the binary..."
	@go build -o $(BUILD_OUT_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

deps:
	@go mod download
	@go mod verify

godoc:
	@echo "Generating godoc..."
	@godoc -http=:6060 -play

help:
	@echo "Makefile for Mind Palace <Go>"
	@echo "Usage: "
	@echo "  make [target]"
	@echo
	@echo "Targets:"
	@echo "  build       Build the binary"
	@echo "  godoc       Generate godoc"
	@echo "  help        Display this help message"
