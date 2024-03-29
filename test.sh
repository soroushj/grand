#!/usr/bin/bash

set -e

TMP_DIR="$(mktemp -d)"

BIN="$TMP_DIR/grand"
go build -cover -o "$BIN" .

GOCOVERDIR="$TMP_DIR/cover"
mkdir "$GOCOVERDIR"

BIN="$BIN" GOCOVERDIR="$GOCOVERDIR" go test -v -count=1 .
go tool covdata percent -i "$GOCOVERDIR"

rm -rf "$TMP_DIR"
