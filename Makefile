PROJECT = $(shell basename $(CURDIR))
PATH_TO_PPROF ?= /var/folders/...mem.pprof
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

graphviz:
	brew install graphviz

# Requires make graphviz
# update PATH_TO_PPROF with output
pprof-pdf:
	go tool pprof --pdf build/diff-checker $(PATH_TO_PPROF) > pprof.pdf

.PHONY: default test dep build local-build run
