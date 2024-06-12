#!/bin/bash

COVER_DIR=coverage
GO_COVER_DIR="$COVER_DIR/go"
PY_COVER_DIR="$COVER_DIR/py"

COVER_GO=cover.go

mkdir -p "$GO_COVER_DIR"
mkdir -p "$PY_COVER_DIR"

cover:go() {
	go test -v -coverprofile "$GO_COVER_DIR/$COVER_GO.out" ./cli/... ./pkg/models/...
	go tool cover -html "$GO_COVER_DIR/$COVER_GO.out" -o "$GO_COVER_DIR/$COVER_GO.html"

	if [ "$1" = "open" ]; then
		open "$GO_COVER_DIR/$COVER_GO.html"
	fi
}

cover:py() {
	poetry run pytest --cov=pkg/rpc/server --cov-report=html:"$PY_COVER_DIR"

	if [ "$1" = "open" ]; then
		open "$PY_COVER_DIR/index.html"
	fi
}

cover:all() {
	cover:go $1
	cover:py $1
}

case $1 in
go)
	cover:go $2
	;;
py)
	cover:py $2
	;;
all)
	cover:all $2
	;;
*)
	echo "please use 'bash cover.sh go|py|all [open]'"
	;;
esac
