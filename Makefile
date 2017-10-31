BIN := bookmarks
PKG := github.com/nrocco/bookmarks/cmd/bookmarks
VERSION := $(shell git describe --tags --always --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v ${PKG}/vendor/)
GO_FILES := $(shell find * -type d -name vendor -prune -or -name '*.go' -type f | grep -v vendor)

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

.DEFAULT_GOAL: build

build/$(BIN)-$(GOOS)-$(GOARCH): $(GO_FILES)
	mkdir -p build
	go generate github.com/nrocco/bookmarks/...
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o $@ ${PKG}

.PHONY: lint
lint:
	@for file in ${GO_FILES}; do golint $${file}; done

.PHONY: vet
vet:
	@go vet ${PKG_LIST}

.PHONY: test
test:
	@go test ${PKG_LIST}

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: clean
clean:
	rm -rf build

.PHONY: build
build: build/$(BIN)-$(GOOS)-$(GOARCH)

.PHONY: container
container: version
	docker build -t "nrocco/bookmarks:latest" .

.PHONY: push
push: container
	docker push "nrocco/bookmarks:latest"

bindata/favicon.ico:
	convert bindata/apple-touch-icon.png -define icon:auto-resize=64,48,32,16 bindata/favicon.ico
