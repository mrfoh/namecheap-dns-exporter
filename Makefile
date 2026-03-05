BINARY_NAME=ncdns-exporter

.PHONY: build

.DEFAULT_GOAL := build

prepare:
	@echo "Preparing build environment..."
	@mkdir -p bin

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

build: prepare
	@echo "Building for local environment..."
	@CGO_ENABLED=0 go build -o bin/$(BINARY_NAME) cmd/main.go