.PHONY: all build start-grpc-server stop-grpc-server deps deps-go deps-py dev-deps deps-llama rpc clean-rpc rpc-py clean-rpc-py test-py test-go test test-e2e cover db vdb graph godoc migrate

BUILD_OUT_DIR=bin
BINARY_NAME=mind-palace
SOURCE_DIR=.

.EXPORT_ALL_VARIABLES:
PROJECT_ROOT=$(shell pwd)
MP_ENV ?= dev# prod,test,dev
LOG_LEVEL ?= 0

all: build

build: deps rpc
	@echo "> Building the binary..."
	@go mod tidy
	@go build -o $(BUILD_OUT_DIR)/$(BINARY_NAME) $(SOURCE_DIR)

dirs:
	@mkdir -p logs

start-grpc-server: dirs
	@echo "> Starting gRPC server with environment $(MP_ENV)"
	@MP_ENV=$(MP_ENV) poetry run python pkg/rpc/server/server.py &

stop-grpc-server:
	@ps aux | grep "python pkg/rpc/server/server.py" | grep -v grep | awk '{print $$2}' | xargs kill

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
	@rm -rf ./pkg/rpc/client/gen ./pkg/rpc/server/gen

rpc-py:
	@echo "> Compiling '.proto' files..."
	@bash ./scripts/pb-compiler-py.sh

clean-rpc-py:
	@echo "> Removing compiled '.proto' python files..."
	@rm -rf ./pkg/rpc/server/gen

# -s don't capture stdout
# -k <test_name>
test-py:
	MP_ENV=test poetry run pytest $(ARGS)

# -count=1 ignores caching
test-go:
	MP_ENV=test LOG_LEVEL=$(LOG_LEVEL) go test -v $(shell go list ./pkg/... ./cli/...) $(ARGS)

test: test-go test-py
	@echo "> Done"

# locally
test-e2e: start-grpc-server
	@$(MAKE) MP_ENV=test db ARGS=start
	@echo "> Running e2e tests..."
	-@ARGS=$${ARGS:="-count=1 -run '^TestE2ETestSuite'"}; \
		$(MAKE) MP_ENV=test test-go ARGS="$$ARGS"
	@echo "> Stopping grpc server"
	@$(MAKE) MP_ENV=test stop-grpc-server
	@sleep 1 # to avoid connection peer timeout
	@echo "> Dropping database"
	@$(MAKE) MP_ENV=test db ARGS=drop
	@echo "> Stopping database"
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
		google,github.com/google,github.com/lib,github.com/joho,github.com/spf13,github.com/rs/zerolog,gopkg \
		-stoponerror=false \
		-s . | dot -Tpng -o godepgraph.png
	@eog godepgraph.png

godoc:
	@echo "> Generating godoc..."
	@godoc -http=:6060 -play
