BUILD_DIR=bin
VERSION=0.1.0
NAME=gioterm
TAGS=nox11
INSTALL_PREFIX=/usr/local

.PHONY: all build clean run fmt install

all: build

install: build 
	mkdir -p $(INSTALL_PREFIX)/bin
	cp $(BUILD_DIR)/$(NAME) $(INSTALL_PREFIX)/bin/$(NAME)

build: fmt
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(NAME) -tags $(TAGS) .

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(NAME)

fmt:
	go fmt
