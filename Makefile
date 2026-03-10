BUILD_DIR := $(CURDIR)/build
APP_BIN_DIR := $(BUILD_DIR)/app-bins
PLUGIN_DIR := $(BUILD_DIR)/plugins

APP := ./cmd/linter
PLUGIN_SRC := ./plugin
PLUGIN_SO := $(PLUGIN_DIR)/linter.so

LINUX_BUILD_FLAGS := -trimpath -ldflags='-w -s -linkmode external -extldflags "-fno-PIC -static"'
PLUGIN_FLAGS := -buildmode=plugin -trimpath -ldflags='-w -s'

.PHONY: pre-commit
pre-commit:
	@pre-commit run --all-files

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR)

.PHONY: test
test:
	@go test -v ./...

# golangci-lint plugin: build, run

.PHONY: build-plugin
build-plugin:
	@echo "building plugin..."
	@mkdir -p $(PLUGIN_DIR)
	go build $(PLUGIN_FLAGS) -o $(PLUGIN_SO) $(PLUGIN_SRC)

.PHONY: test-plugin
test-plugin: build-plugin
	@echo "running golangci-lint with plugin..."
	@golangci-lint run --config .golangci.yml ./...

# non-plugin builds: native, linux-arm64, linux-amd64

.PHONY: build
build:
	@echo "building native..."
	@mkdir -p $(APP_BIN_DIR)/native
	go build -o $(APP_BIN_DIR)/native/linter $(APP)

.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "building for Linux ARM64 with musl (static)..."
	@mkdir -p $(APP_BIN_DIR)/linux_arm64
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-musl-gcc CXX=aarch64-linux-musl-g++ go build $(LINUX_BUILD_FLAGS) -v -o $(APP_BIN_DIR)/linux_arm64/linter $(APP)

.PHONY: build-linux-amd64
build-linux-amd64:
	@echo "building for Linux AMD64 with musl (static)..."
	@mkdir -p $(APP_BIN_DIR)/linux_amd64
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build $(LINUX_BUILD_FLAGS) -v -o $(APP_BIN_DIR)/linux_amd64/linter $(APP)
