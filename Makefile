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
	go generate pkg/server/*.go
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -x -v -a -o $@ ${PKG}

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
container:
	docker build --pull -t "nrocco/bookmarks" .

bindata/favicon.ico:
	convert bindata/apple-touch-icon.png -define icon:auto-resize=64,48,32,16 bindata/favicon.ico
