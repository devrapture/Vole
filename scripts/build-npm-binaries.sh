#!/usr/bin/env bash
set -euo pipefail

mkdir -p dist

GOOS=darwin GOARCH=arm64 go build -o dist/vole-clean-darwin-arm64 main.go
GOOS=darwin GOARCH=amd64 go build -o dist/vole-clean-darwin-amd64 main.go

GOOS=linux GOARCH=arm64 go build -o dist/vole-clean-linux-arm64 main.go
GOOS=linux GOARCH=amd64 go build -o dist/vole-clean-linux-amd64 main.go

GOOS=windows GOARCH=amd64 go build -o dist/vole-clean-windows-amd64.exe main.go

chmod +x dist/vole-clean-*
chmod +x npm/vole-clean.js
