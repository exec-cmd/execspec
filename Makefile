.DEFAULT_GOAL := build

.PHONY: fmt vet test build

fmt: 
	go fmt ./...

vet: fmt
	go vet ./...

test: vet
	go test ./...

build: test
	mkdir -p build/bin
	go build -o build/bin/xspec ./src/xspec
