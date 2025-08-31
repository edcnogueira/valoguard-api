# Makefile for Valoguard CLI

# Variables
BINARY ?= bin/valoguard
PACKAGE ?= ./cmd/cli
API_KEY ?=
LDFLAGS ?=

# If API_KEY provided, append ldflag to embed it
ifneq ($(API_KEY),)
	LDFLAGS += -X github.com/edcnogueira/valoguard-api/cmd/cli/transport/player.defaultAPIKey=$(API_KEY)
endif

.PHONY: help build-cli clean

help:
	@echo "Valoguard Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  make build-cli API_KEY=HDEV-xxxx   Build CLI with embedded Henrik API key"
	@echo "  make clean                         Remove build artifacts"
	@echo ""
	@echo "Variables:"
	@echo "  BINARY  (default: bin/valoguard)"
	@echo "  PACKAGE (default: ./cmd/cli)"

build-cli:
	@test -n "$(API_KEY)" || (echo "Error: API_KEY not defined. Use: make build-cli API_KEY=HDEV-xxxx"; exit 1)
	@mkdir -p $(dir $(BINARY))
	@echo "Building CLI with embedded key..."
	@go build -o $(BINARY) -ldflags "$(LDFLAGS)" $(PACKAGE)
	@echo "Done: $(BINARY)"

clean:
	@rm -rf bin
	@echo "Cleaned"
