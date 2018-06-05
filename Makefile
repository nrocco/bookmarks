BIN := bookmarks
PKG := github.com/nrocco/bookmarks
CONTAINER := nrocco/bookmarks
VERSION := $(shell git describe --tags --always --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v ${PKG}/vendor/)
GO_FILES := $(shell git ls-files '*.go')

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

.DEFAULT_GOAL: build

build/${BIN}-$(GOOS)-$(GOARCH): frontend $(GO_FILES)
	mkdir -p build
	go generate ${PKG}/...
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -v -o $@ -ldflags "-X main.Version=${VERSION}" ${PKG}/cmd/${BIN}

.PHONY: lint
lint:
	golint -set_exit_status ${PKG_LIST}

.PHONY: vet
vet:
	go vet -v ${PKG_LIST}

.PHONY: test
test:
	go test -short ${PKG_LIST}

.PHONY: coverage
coverage:
	mkdir -p coverage && rm -rf coverage/*
	for package in ${PKG_LIST}; do go test -covermode=count -coverprofile "coverage/$${package##*/}.cov" "$$package"; done
	echo mode: count > coverage/coverage.cov
	tail -q -n +2 coverage/*.cov >> coverage/coverage.cov
	go tool cover -func=coverage/coverage.cov

.PHONY: version
version:
	@echo ${VERSION}

.PHONY: clean
clean:
	rm -rf build

.PHONY: build
build: build/${BIN}-${GOOS}-${GOARCH}

.PHONY: container
container: version
	docker build --build-arg "VERSION=${VERSION}" -t "${CONTAINER}:${VERSION}" .

.PHONY: push
push: container
	docker push "${CONTAINER}:${VERSION}"

.PHONY: frontend
frontend:
	$(MAKE) -C frontend build

.PHONY: server
server:
	modd
