.DEFAULT_GOAL := build

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: fmt vet test build

fmt: 
	go fmt ./...

vet: fmt
	go vet ./...

test: vet
	go test ./...

build: test
	mkdir -p build/bin
	go build -ldflags "-X github.com/exec-cmd/execspec/src/cli.version=$(VERSION)" -o build/bin/xspec ./src/xspec
