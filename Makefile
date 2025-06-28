# We're going to be using this Makefile as a sort of task runner, for all sorts of operations in this project

# first we'll grab the current version from our ENV VAR (added by our CI) - see here: https://github.com/marketplace/actions/version-increment
BINARY_NAME := otdfctl
CURR_VERSION := ${SEM_VER}
COMMIT_SHA := ${COMMIT_SHA}
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')

GO_MOD_LINE = $(shell head -n 1 go.mod | cut -c 8-)
GO_MOD_NAME = $(word 1,$(subst /, ,$(GO_MOD_LINE)))
APP_CFG = $(GO_MOD_LINE)/pkg/config

GO_BUILD_FLAGS=-ldflags " \
	-X $(APP_CFG).Version=${CURR_VERSION} \
	-X $(APP_CFG).CommitSha=${COMMIT_SHA} \
	-X $(APP_CFG).BuildTime=${BUILD_TIME} \
"
GO_BUILD_PREFIX=$(TARGET_DIR)/$(BINARY_NAME)-${CURR_VERSION}

# If commit sha is not available try git
ifndef COMMIT_SHA
	COMMIT_SHA := $(shell git rev-parse HEAD)
endif

# If current version is not available try git
ifndef CURR_VERSION
	CURR_VERSION := $(shell git describe --tags --always)
endif

# Default target executed when no arguments are given to make.
# NOTE: .PHONY is used to indicate that the target is not a file (e.g. there is no file called 'build-darwin-amd64', instead the .PHONY directive tells make that the proceeding target is a command to be executed, not a file to be generated)
.PHONY: all
all: run
.DEFAULT_GOAL := run

# Target directory for compiled binaries
TARGET_DIR=target

# Output directory for the zipped artifacts
OUTPUT_DIR=output

# Build commands for each platform (extra hyphen used in windows to avoid issues with the .exe extension)
PLATFORMS := \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	linux-arm \
	linux-arm64 \
	windows-amd64-.exe \
	windows-arm-.exe \
	windows-arm64-.exe

build: test clean $(addprefix build-,$(PLATFORMS)) zip-builds verify-checksums

build-%:
	GOOS=$(word 1,$(subst -, ,$*)) \
	GOARCH=$(word 2,$(subst -, ,$*)) \
	go build $(GO_BUILD_FLAGS) \
		-o $(GO_BUILD_PREFIX)-$(word 1,$(subst -, ,$*))-$(word 2,$(subst -, ,$*))$(word 3,$(subst -, ,$*))

zip-builds:
	./.github/scripts/zip-builds.sh $(BINARY_NAME)-$(CURR_VERSION) $(TARGET_DIR) $(OUTPUT_DIR)

verify-checksums:
	./.github/scripts/verify-checksums.sh $(OUTPUT_DIR) $(BINARY_NAME)-$(CURR_VERSION)_checksums.txt 

# Target for running the project (adjust as necessary for your project)
.PHONY: run
run:
	go run .

# Target for testing the project
.PHONY: test
test:
	go test -v ./...

.PHONY: build-test
build-test:
	go build \
		-ldflags "\
			-X $(APP_CFG).TestMode=true \
			-X $(APP_CFG).Version=${CURR_VERSION}-testbuild \
			-X $(APP_CFG).CommitSha=${COMMIT_SHA} \
			-X $(APP_CFG).BuildTime=${BUILD_TIME} \
		" \
		-o $(BINARY_NAME)_testbuild

.PHONY: test-bats
test-bats: build-test
	./e2e/resize_terminal.sh && bats ./e2e

# Target for cleaning up the target directory
.PHONY: clean
clean:
	rm -rf $(TARGET_DIR)

# Target for generating CLI commands from documentation
# NOTE: Generated files require manual implementation of handler interfaces
# and should be committed to the repository after implementation
.PHONY: generate
generate:
	go run github.com/jrschumacher/adder/cmd/adder@v0.1.1 generate -o cmd/generated

# Target for cleaning generated files
.PHONY: clean-generated
clean-generated:
	rm -rf cmd/generated/

# Target for regenerating (clean + generate)
# WARNING: This will remove existing generated files - ensure handlers are implemented elsewhere
.PHONY: regenerate
regenerate: clean-generated generate

# Target for checking required tools are available
.PHONY: toolcheck
toolcheck:
	@echo "Checking required tools..."
	@command -v go >/dev/null 2>&1 || { echo >&2 "go is required but not installed. Visit https://golang.org/dl/"; exit 1; }
	@echo "✓ go is available"
	@go version 2>/dev/null | grep -q "go1\." || { echo >&2 "go version check failed"; exit 1; }
	@echo "✓ go version is compatible"
	@go run github.com/jrschumacher/adder/cmd/adder@v0.1.1 version >/dev/null 2>&1 || { echo >&2 "adder tool check failed. Run 'go mod download' to ensure dependencies are available."; exit 1; }
	@echo "✓ adder is available"
	@command -v bats >/dev/null 2>&1 || echo "⚠ bats is not installed (required for e2e tests). Install with: brew install bats-core"
	@echo "All required tools are available!"
