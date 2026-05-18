BINARY     := ollama-oneshot
MODULE     := github.com/kasjfulk/ollama-oneshot
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.1.0")
BUILD_DIR  := dist
INSTALL_DIR := $(HOME)/.local/bin

LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: all build install uninstall clean test lint run help

all: build

## build: compile binary to ./dist/
build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) .
	@echo "Built $(BUILD_DIR)/$(BINARY)"

## install: build and install to ~/.local/bin
install: build
	@mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed to $(INSTALL_DIR)/$(BINARY)"

## uninstall: remove installed binary
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY)
	@echo "Removed $(INSTALL_DIR)/$(BINARY)"

## clean: remove build artifacts
clean:
	rm -rf $(BUILD_DIR)

## test: run all tests
test:
	go test ./...

## lint: run go vet
lint:
	go vet ./...

## run: build and run with ARGS (e.g. make run ARGS="run --prompt 'hello' --dry-run")
run: build
	$(BUILD_DIR)/$(BINARY) $(ARGS)

## deps: tidy and download dependencies
deps:
	go mod tidy
	go mod download

## help: show this help
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
