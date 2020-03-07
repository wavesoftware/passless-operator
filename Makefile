PROJECT_DIR   = $(shell readlink -f .)
BUILD_DIR     = "$(PROJECT_DIR)/build/_output"
MANAGER_DIR   = "$(PROJECT_DIR)/cmd/manager"
MANAGER_BIN   = "$(BUILD_DIR)/bin/passless-operator"
VERSION       = $(shell git describe --always --dirty)

GO           ?= go
RICHGO       ?= rich$(GO)

.PHONY: default
default: binary

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
	revive -config revive.toml -formatter stylish ./...
	staticcheck -f stylish ./...

.PHONY: test
test: builddir check
	@echo " ✔️ Testing"
	$(RICHGO) test -v \
		-covermode=count -coverprofile=$(BUILD_DIR)/coverage.out \
		-ldflags="-X github.com/wavesoftware/passless-operator/version.Version=$(VERSION)" \
		./...

.PHONY: binary
binary: builddir test
	@echo " 🔨 Building binary"
	$(RICHGO) build \
		-ldflags="-X github.com/wavesoftware/passless-operator/version.Version=$(VERSION)" \
		-o $(MANAGER_BIN) $(MANAGER_DIR)

.PHONy: image-sole
image-sole:
	@echo " 🔨 Building image"
	operator-sdk build quay.io/wavesoftware/passless-operator

.PHONY: image
image: binary image-sole
