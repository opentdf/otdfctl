# We're going to be using this Makefile as a sort of task runner, for all sorts of operations in this project

# first we'll grab the current version from our ENV VAR (added by our CI) - see here: https://github.com/marketplace/actions/version-increment
CURR_VERSION := ${SEM_VER}

# Default target executed when no arguments are given to make.
# NOTE: .PHONY is used to indicate that the target is not a file (e.g. there is no file called 'build-darwin-amd64', instead the .PHONY directive tells make that the proceeding target is a command to be executed, not a file to be generated)
.PHONY: all
all: run
.DEFAULT_GOAL := run




# Binary name: Change this to your actual binary name
BINARY_NAME=${BIN_NAME}


# Target directory for compiled binaries
TARGET_DIR=target

# Output directory for the zipped artifacts
OUTPUT_DIR=output

# Targets for building the project for different platforms
.PHONY: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm build-linux-arm64 build-windows-amd64 build-windows-arm build-windows-arm64
build: clean build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm build-linux-arm64 build-windows-amd64 build-windows-arm build-windows-arm64

# Build commands for each platform
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-darwin-amd64 .

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-darwin-arm64 .

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-linux-amd64 .

build-linux-arm:
	GOOS=linux GOARCH=arm go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-linux-arm .

build-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-linux-arm64 .

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-windows-amd64.exe .

build-windows-arm:
	GOOS=windows GOARCH=arm go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-windows-arm.exe .

build-windows-arm64:
	GOOS=windows GOARCH=arm64 go build -o $(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}-windows-arm64.exe .

# Target for running the project (adjust as necessary for your project)
.PHONY: run
run: build
	go run .

# Target for testing the project
.PHONY: test
test: build
	go test -v ./...

# Target for cleaning up the target directory
.PHONY: clean
clean:
	rm -rf $(TARGET_DIR)

# Script for zipping up the compiled binaries
.PHONY: zip-builds
zip-builds:
	./.github/scripts/zip-builds.sh $(BINARY_NAME)-$(CURR_VERSION) $(TARGET_DIR) $(OUTPUT_DIR)

# Script for verifying the checksums
.PHONY: verify-checksums
verify-checksums:
	.github/scripts/verify-checksums.sh $(OUTPUT_DIR) $(BINARY_NAME)-$(CURR_VERSION)_checksums.txt 