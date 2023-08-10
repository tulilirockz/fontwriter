TARGET := fontwriter
BUILD_DIR := build
DIST_FOLDER := dist
DATE_FORMAT := "+%H-%d-%m-%Y"
GIT_CURR_HASH := $(shell git rev-parse --short HEAD)

ifdef OS	
	GO_TARGET := windows
	TARGET_EXTENSION := exe
	PKG_FMT := zip
	CURRENT_OS := windows
	PACKAGING_COMMAND := zip -r
else
   	ifeq ($(shell uname -s), Linux)
		GO_TARGET := linux
		TARGET_EXTENSION := elf
		PKG_FMT := tar.gz
		CURRENT_OS := linux
		PACKAGING_COMMAND := tar czvf
	endif
endif

all: package

build:
	mkdir -p $(BUILD_DIR)
	GOOS=$(GO_TARGET) GO_ARCH=amd64 go build -o $(BUILD_DIR)/$(TARGET).$(TARGET_EXTENSION)
	cp config.toml build

package: clean build
	mkdir -p $(DIST_FOLDER)
	$(PACKAGING_COMMAND) $(DIST_FOLDER)/$(BUILD_DIR)-$(CURRENT_OS)-$(GIT_CURR_HASH).$(PKG_FMT) $(BUILD_DIR)

.PHONY: clean
clean:
	rm -rf $(DIST_FOLDER) $(BUILD_DIR) out/