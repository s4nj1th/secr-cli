# Makefile for secr-cli
BINARY := secr-cli
INSTALL_PATH := /usr/local/bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X secr-cli/cmd.version=$(VERSION)"

.PHONY: build install uninstall test clean help

build:
	go build $(LDFLAGS) -o $(BINARY) .

install: build
	@echo "Installing $(BINARY) to $(INSTALL_PATH)"
	@mkdir -p $(INSTALL_PATH)
	@mv $(BINARY) $(INSTALL_PATH)
	@chmod +x $(INSTALL_PATH)/$(BINARY)
	@echo "Installed. Run '$(BINARY) --help' to get started."

uninstall:
	@echo "Removing $(INSTALL_PATH)/$(BINARY)"
	@rm -f $(INSTALL_PATH)/$(BINARY)
	@echo "Uninstalled."

test:
	go test ./... -v

clean:
	rm -f $(BINARY)

help:
	@echo "Available targets:"
	@echo "  build      - Compile the binary"
	@echo "  install    - Install to $(INSTALL_PATH)"
	@echo "  uninstall  - Remove from system"
	@echo "  test       - Run test suite"
	@echo "  clean      - Remove build artifacts"