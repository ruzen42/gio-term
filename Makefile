BUILD_DIR=output
VERSION=0.1.0
NAME=main

.PHONY: all build clean run fmt install

all: build

install: build 
	mkdir -p /usr/local/bin
	cp $(BUILD_DIR)/$(NAME) /usr/local/bin/$(NAME)

build: fmt
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(NAME) .

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(NAME)

fmt:
	go fmt
