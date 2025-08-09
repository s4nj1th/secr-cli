BINARY=secr-cli
INSTALL_PATH=/usr/local/bin

.PHONY: build install

build:
	go build -o $(BINARY) ./cmd

install: build
	cp $(BINARY) $(INSTALL_PATH)/
	chmod +x $(INSTALL_PATH)/$(BINARY)
