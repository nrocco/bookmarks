.DEFAULT_GOAL := build

BIN ?= bookmarks
REPO ?= nrocco/bookmarks
PKG ?= github.com/nrocco/bookmarks
DOCKER_IMAGE = nrocco/bookmarks

CGO_ENABLED ?= 1
BUILD_GOOS ?= $(shell go env GOOS)
BUILD_GOARCH ?= $(shell go env GOARCH)
BUILD_VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_COMMIT ?= $(shell git describe --always --dirty)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILD_NAME ?= $(BIN)-$(BUILD_VERSION)-$(BUILD_GOOS)-$(BUILD_GOARCH)


build: lint test dist/$(BUILD_NAME)/bin/$(BIN)


archive: dist/$(BUILD_NAME).tar.gz
	tar tf "$<"


dist/$(BUILD_NAME)/bin/$(BIN):
	mkdir -p "$(@D)"
	env GOOS=$(BUILD_GOOS) GOARCH=$(BUILD_GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build \
		-v \
		-o "$@" \
		-ldflags "-X ${PKG}/cmd.version=${BUILD_VERSION} -X ${PKG}/cmd.commit=${BUILD_COMMIT} -X ${PKG}/cmd.date=${BUILD_DATE}"


dist/$(BUILD_NAME).tar.gz: dist/$(BUILD_NAME)/bin/$(BIN) dist/completion.zsh bin/* LICENSE README.md
	mkdir -p "dist/$(BUILD_NAME)"
	cp bin/* "dist/$(BUILD_NAME)/bin"
	cp LICENSE README.md dist/completion.zsh "dist/$(BUILD_NAME)"
	tar czf "dist/$(BUILD_NAME).tar.gz" -C dist/ "$(BUILD_NAME)"


server/server.pb.go: server/server.proto
	protoc -I server/ server/server.proto --go_out=plugins=grpc:server


dist/completion.zsh:
	$(MAKE) build BUILD_GOARCH=$(shell go env GOARCH) BUILD_GOOS=$(shell go env GOOS)
	dist/$(BIN)-$(BUILD_VERSION)-$(shell go env GOOS)-$(shell go env GOARCH)/bin/$(BIN) completion > "$@"


.PHONY: clear
clear:
	rm -rf dist


.PHONY: build-all
build-all:
	$(MAKE) build BUILD_GOARCH=amd64 BUILD_GOOS=darwin
	$(MAKE) build BUILD_GOARCH=amd64 BUILD_GOOS=freebsd
	$(MAKE) build BUILD_GOARCH=amd64 BUILD_GOOS=linux


.PHONY: archive-all
archive-all:
	$(MAKE) archive BUILD_GOARCH=amd64 BUILD_GOOS=darwin
	$(MAKE) archive BUILD_GOARCH=amd64 BUILD_GOOS=freebsd
	$(MAKE) archive BUILD_GOARCH=amd64 BUILD_GOOS=linux


.PHONY: release
release: archive-all
	sha256sum dist/*.tar.gz > dist/checksums.txt
	tools/release-to-github.py $(REPO) $(BUILD_VERSION) dist/checksums.txt dist/*.tar.gz


.PHONY: lint
lint:
	git ls-files | xargs misspell -error
	golint -set_exit_status ./...
	go vet -v ./...
	errcheck -blank -asserts ./...

.PHONY: test
test:
	go test -v -short ./...

.PHONY: coverage
coverage:
	mkdir -p coverage
	go test -covermode=count -coverprofile=coverage/coverage.out ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

.PHONY: test
test:
	go test ./...

.PHONY: container
container:
	docker build \
		--build-arg "VERSION=$(BUILD_VERSION)" \
		--build-arg "COMMIT=$(BUILD_COMMIT)" \
		--build-arg "DATE=$(BUILD_DATE)" \
		--tag "$(DOCKER_IMAGE):latest" \
		.

.PHONY: push
push:
	docker push "$(DOCKER_IMAGE):latest"
