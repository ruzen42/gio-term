BUILD_DIR=bin
VERSION=0.1.0
NAME=gioterm
TAGS=nox11

.PHONY: all build clean run fmt install

all: build

install: build 
	mkdir -p /usr/local/bin
	cp $(BUILD_DIR)/$(NAME) /usr/local/bin/$(NAME)

build: fmt
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(NAME) -tags $(TAGS) .

clean:
	rm -rf $(BUILD_DIR)

run: build
	./$(BUILD_DIR)/$(NAME)

fmt:
	go fmt
