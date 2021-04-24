DOCKER_IMAGE = nrocco/bookmarks
DOCKER_IMAGE_VERSION = latest
BUILD_VERSION ?= $(shell git describe --tags --always --dirty)
BUILD_COMMIT ?= $(shell git describe --always --dirty)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

.PHONY: help
help:
	@echo 'make build-all dist/bookmarks-amd64-freebsd dist/bookmarks-amd64-darwin dist/bookmarks-amd64-linux clear container push'

.PHONY: build-all
build-all: \
	dist/bookmarks-amd64-freebsd \
	dist/bookmarks-amd64-darwin \
	dist/bookmarks-amd64-linux

.PHONY: dist/bookmarks-amd64-freebsd
dist/bookmarks-amd64-freebsd:
	mkdir -p dist/bookmarks-amd64-freebsd
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--target bin \
		--platform freebsd/amd64 \
		--output dist/bookmarks-amd64-freebsd \
		.

.PHONY: dist/bookmarks-amd64-darwin
dist/bookmarks-amd64-darwin:
	mkdir -p dist/bookmarks-amd64-darwin
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--target bin \
		--platform darwin/amd64 \
		--output dist/bookmarks-amd64-darwin \
		.

.PHONY: dist/bookmarks-amd64-linux
dist/bookmarks-amd64-linux:
	mkdir -p dist/bookmarks-amd64-linux
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--target bin \
		--platform linux/amd64 \
		--output dist/bookmarks-amd64-linux \
		.

.PHONY: clear
clear:
	rm -rf dist

.PHONY: container
container:
	docker image build --pull \
		--build-arg "BUILD_VERSION=$(BUILD_VERSION)" \
		--build-arg "BUILD_COMMIT=$(BUILD_COMMIT)" \
		--build-arg "BUILD_DATE=$(BUILD_DATE)" \
		--tag "$(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)" \
		.

.PHONY: push
push: container
	docker image push "$(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)"
