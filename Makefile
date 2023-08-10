TARGET := fontwriter
BUILD_DIR := build
DIST_FOLDER := dist

ifdef OS
	GO_TARGET := windows
	TARGET_EXTENSION := exe
else
   ifeq ($(shell uname), Linux)
      GO_TARGET := linux
	  TARGET_EXTENSION := elf
   endif
endif

all: build

build:
	mkdir -p $(BUILD_DIR)
	GOOS=$(GO_TARGET) GO_ARCH=amd64 go build -o $(BUILD_DIR)/$(TARGET).$(TARGET_EXTENSION)
	cp config.toml build

package: clean build
	mkdir -p $(DIST_FOLDER)
	tar czvf $(DIST_FOLDER)/$(BUILD_DIR).tar.gz $(BUILD_DIR)

build_win:
	OS="windows" make build

package_win: clean build_win
	mkdir -p $(DIST_FOLDER)
	zip -r $(DIST_FOLDER)/$(BUILD_DIR).zip $(BUILD_DIR)

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(DIST_FOLDER) $(BUILD_DIR) out/