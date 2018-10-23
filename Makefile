PROJECT = $(shell basename $(CURDIR))

default: test

test:
	go test -v -race ./src/...

dep:
	go get -u github.com/golang/dep/cmd/dep

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/$(PROJECT) ./cmd/$(PROJECT);

local-build:
	@go build -o build/$(PROJECT) ./cmd/$(PROJECT);

run: local-build
	./build/$(PROJECT)

.PHONY: default test dep build local-build run
