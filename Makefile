TARGET := fontwriter
BUILD_DIR := build

ifdef OS
	GO_TARGET := windows
else
   ifeq ($(shell uname), Linux)
      GO_TARGET := linux
   endif
endif

all: build

build:
	mkdir -p $(BUILD_DIR)
	GOOS=$(GO_TARGET) GO_ARCH=amd64 go build -o $(BUILD_DIR)/$(TARGET)
	cp config.toml build

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf out/