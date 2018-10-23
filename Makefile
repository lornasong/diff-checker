PROJECT = $(shell basename $(CURDIR))

.PHONY: default test build local-build run

default: test

test:
	go test -v -race ./src/...

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$(PROJECT) ./cmd/$(PROJECT);

local-build:
	@go build -o build/$(PROJECT) ./cmd/$(PROJECT);

run: local-build
	./build/$(PROJECT)
