.PHONY: default build clean

.DEFAULT_GOAL := default

default: clean generate build

# Detect architecture and OS
GOARCH ?= $(shell go env GOARCH)
GOOS ?= $(shell go env GOOS)

# Check that OS is linux
ifneq ($(GOOS),linux)
    $(error Unsupported OS: $(GOOS). Only linux is supported)
endif

# Check that architecture is supported
ifneq ($(GOARCH),amd64)
ifneq ($(GOARCH),arm64)
    $(error Unsupported architecture: $(GOARCH). Only amd64 and arm64 are supported)
endif
endif

# Output directory and binary name
BIN_DIR := bin
BINARY_NAME := no-redirects-$(GOOS)-$(GOARCH).so

# ---------- build ----------

build:
	@echo "Building for GOARCH=$(GOARCH) -> $(BINARY_NAME)"
	@mkdir -p $(BIN_DIR)
	CGO_ENABLED=1 GOOS=linux GOARCH=$(GOARCH) go build -buildmode=plugin -o $(BIN_DIR)/$(BINARY_NAME) .

# ---------- clean ----------

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
