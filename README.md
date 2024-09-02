# Mind Palace

## Stack

- Golang >= 1.22
- Python >= 3.9

### Additionally

- `protoc` >= 3.21
- `protoc-gen-go` 1.28
- `protoc-gen-go-grpc` 1.2

## Installation

```bash
make deps
```

## Helpful

```bash
# disable test cache
make test-go ARGS="-count=1"
# run test names containing TestSuite
make test-go ARGS="-run TestSuite"
# regex also works
make test-go ARGS='-run "^TestSuite/SomeTest.*$"'
```
