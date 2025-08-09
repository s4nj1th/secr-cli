BINARY=secr-cli
INSTALL_PATH=/usr/local/bin

.PHONY: build install

build:
	go build -o $(BINARY) ./main.go

install: build
	cp $(BINARY) $(INSTALL_PATH)/
	chmod +x $(INSTALL_PATH)/$(BINARY)
