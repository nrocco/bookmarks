BIN := bookmarks
PKG := github.com/nrocco/bookmarks/cmd/bookmarks
VERSION := $(shell git describe --tags --always --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v ${PKG}/vendor/)
GO_FILES := $(shell find * -type d -name vendor -prune -or -name '*.go' -type f | grep -v vendor)

PREFIX = /usr/local

.DEFAULT_GOAL: build/$(BIN)

build/$(BIN): $(GO_FILES)
	CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo -ldflags "-d -s -w -X ${PKG}/cmd.Version=${VERSION}" -o build/${BIN} ${PKG}

.PHONY: deps
deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

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

public/favicon.ico:
	convert public/apple-touch-icon.png -define icon:auto-resize=64,48,32,16 public/favicon.ico

test-server:
	go run cmd/bookmarks/main.go --database 'postgres://postgres:secret@localhost/bookmarks?sslmode=disable'
