.PHONY: all build start-palace-grpc-server start-vdb-grpc-server start-log-grpc-server stop-palace-grpc-server stop-vdb-grpc-server stop-log-grpc-server deps deps-go deps-py dev-deps deps-llama rpc clean-rpc test-py test-go test-go-helper test test-e2e cover db vdb graph godoc migrate

BUILD_OUT_DIR=bin
BINARY_NAME=mind-palace
SOURCE_DIR=.

.EXPORT_ALL_VARIABLES:
PROJECT_ROOT=$(shell pwd)
MP_ENV ?= dev# prod,test,dev

# https://stackoverflow.com/a/70663753/14414945
ifneq (,$(wildcard ./.env.$(MP_ENV)))
    include .env.$(MP_ENV)
    export
endif

all: build

build: deps rpc
	@echo "> Building the binary..."
	@go mod tidy
	@go build -o $(BUILD_OUT_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

dirs:
	@mkdir -p logs

start-palace-grpc-server: dirs
	@echo "> Starting Palace gRPC server with environment $(MP_ENV)"
	@poetry run python pkg/rpc/palace/server.py

start-vdb-grpc-server: dirs
	@echo "> Starting VDB gRPC server with environment $(MP_ENV)"
	@poetry run python pkg/rpc/vdb/server.py

start-log-grpc-server: dirs
	@echo "> Starting Log gRPC server with environment $(MP_ENV)"
	@poetry run python pkg/rpc/log/server.py

stop-palace-grpc-server:
	@echo "> Stopping Palace gRPC server with environment $(MP_ENV)"
	@ps aux | grep "python pkg/rpc/palace/server.py" | grep -v grep | awk '{print $$2}' | xargs kill

stop-vdb-grpc-server:
	@echo "> Stopping VDB gRPC server with environment $(MP_ENV)"
	@ps aux | grep "python pkg/rpc/vdb/server.py" | grep -v grep | awk '{print $$2}' | xargs kill

stop-log-grpc-server:
	@echo "> Stopping Log gRPC server with environment $(MP_ENV)"
	@ps aux | grep "python pkg/rpc/log/server.py" | grep -v grep | awk '{print $$2}' | xargs kill

deps: deps-go deps-py

deps-go:
	@echo "> Installing go dependencies..."
	@go mod download
	@go mod verify

deps-py:
	@echo "> Installing python dependencies..."
	@poetry install

dev-deps:
	@go install github.com/kisielk/godepgraph@latest
	@go install golang.org/x/tools/cmd/cover@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
	@poetry add pytest-cov

# export PATH="/usr/local/cuda-12.5/bin:$PATH"
# export CUDA_HOME=/usr/local/cuda-12.5
deps-llama:
	# CMAKE_ARGS="-DLLAMA_CUDA=on" LLAMA_CCACHE=OFF FORCE_CMAKE=1 poetry run pip install llama-index-core==0.10.43 llama-index-llms-llama-cpp==0.1.3 --no-cache-dir --force-reinstall --upgrade
	CMAKE_ARGS="-DLLAMA_CUDA=on" LLAMA_CCACHE=OFF FORCE_CMAKE=1 poetry run pip install llama-cpp-python --no-cache-dir --force-reinstall --upgrade

rpc:
	@echo "> Compiling '.proto' files..."
	@bash ./scripts/pb-compiler.sh

clean-rpc:
	@echo "> Removing compiled '.proto' files..."
	@rm -rf ./pkg/rpc/gen

# -s don't capture stdout
# -k <test_name>
test-py:
	MP_ENV=test poetry run pytest $(ARGS)

test-go:
	@$(MAKE) MP_ENV=test test-go-helper ARGS="-tags=e2e $$ARGS"

# -count=1 ignores caching
test-go-helper:
	MP_ENV=test LOG_LEVEL=$(LOG_LEVEL) go test -v $(shell go list ./pkg/... ./cli/...) $(ARGS)

test:
	@$(MAKE) MP_ENV=test start-log-grpc-server &
	-@$(MAKE) test-go
	-@$(MAKE) test-py
	@$(MAKE) MP_ENV=test stop-log-grpc-server

# locally
test-e2e:
	@$(MAKE) MP_ENV=test db ARGS=start
	@$(MAKE) MP_ENV=test start-log-grpc-server &
	@$(MAKE) MP_ENV=test start-palace-grpc-server &
	@$(MAKE) MP_ENV=test start-vdb-grpc-server &
	@echo "> Running e2e tests..."
	-@ARGS=$${ARGS:="-count=1 -run '^TestE2ETestSuite'"}; \
		$(MAKE) MP_ENV=test test-go-helper ARGS="$$ARGS"
	@$(MAKE) MP_ENV=test stop-palace-grpc-server
	@$(MAKE) MP_ENV=test stop-vdb-grpc-server
	@$(MAKE) MP_ENV=test stop-log-grpc-server
	@sleep 1 # to avoid connection peer timeout
	@$(MAKE) MP_ENV=test db ARGS=drop
	@$(MAKE) MP_ENV=test db ARGS=stop

cover: dev-deps
	MP_ENV=test LOG_LEVEL=5 bash scripts/cover.sh $(ARGS)

db:
	@bash scripts/postgres.sh $(ARGS)

vdb:
	@bash scripts/standalone_embed.sh $(ARGS)

migrate:
	@go run . migrate $(ARGS)

graph:
	@godepgraph -p \
		google,github.com/google,github.com/lib,github.com/joho,github.com/spf13,github.com/rs/zerolog,gopkg,github.com/golang-migrate \
	-stoponerror=false \
		-s . | dot -Tpng -o godepgraph.png
	@eog godepgraph.png

godoc:
	@echo "> Generating godoc..."
	@godoc -http=:6060 -play
