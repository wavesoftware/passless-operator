PROJECT_DIR   = $(shell readlink -f .)
BUILD_DIR     = "$(PROJECT_DIR)/build/_output"
MANAGER_DIR   = "$(PROJECT_DIR)/cmd/manager"
MANAGER_BIN   = "$(BUILD_DIR)/bin/passless-operator"

GO           ?= go
RICHGO       ?= rich$(GO)

.PHONY: default
default: binaries

.PHONY: builddeps
builddeps:
	@GO111MODULE=off $(GO) get github.com/kyoh86/richgo
	@GO111MODULE=off $(GO) get github.com/mgechev/revive
	@GO111MODULE=off $(GO) get honnef.co/go/tools/cmd/staticcheck

.PHONY: builddir
builddir:
	@mkdir -p $(BUILD_DIR)/bin

.PHONY: clean
clean: builddeps
	@echo " 🛁 Cleaning"
	@rm -frv $(BUILD_DIR)

.PHONY: check
check: builddeps
	@echo " 🛂 Checking"
	staticcheck -f stylish ./...
	revive -config revive.toml -formatter stylish ./...

.PHONY: test
test: builddir check
	@echo " ✔️ Testing"
	$(RICHGO) test -v -covermode=count -coverprofile=$(BUILD_DIR)/coverage.out ./...

.PHONY: manager
manager: builddir test
	@echo " 🔨 Building manager"
	$(RICHGO) build -o $(MANAGER_BIN) $(MANAGER_DIR)

.PHONY: binaries
binaries: manager

.PHONY: images
images:
	@echo " 🔨 Building images"
	operator-sdk build quay.io/wavesoftware/passless-operator
