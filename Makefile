# Simple Makefile for local build and install
BINARY := secr-cli
INSTALL_PATH := /usr/local/bin

.PHONY: build install

build:
	go build -o $(BINARY) main.go

install: build
	@echo "Installing $(BINARY) to $(INSTALL_PATH)"
	@mkdir -p $(INSTALL_PATH)
	@mv $(BINARY) $(INSTALL_PATH)
	@chmod +x $(INSTALL_PATH)/$(BINARY)
	@echo "Installed. Run '$(BINARY)' to use."

uninstall:
	@echo "Removing $(INSTALL_PATH)/$(BINARY)"
	@rm -f $(INSTALL_PATH)/$(BINARY)
	@echo "Uninstalled."

help:
	@echo "Available targets:"
	@echo "  build    - Compile the binary"
	@echo "  install  - Install to $(INSTALL_PATH)"
	@echo "  uninstall - Remove from system"