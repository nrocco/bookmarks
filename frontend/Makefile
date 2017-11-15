VERSION := $(shell git describe --tags --always --dirty)

.PHONY: build
build:
	VERSION=${VERSION} npm run build --production --no-progress

.PHONY: lint
lint:
	npm run lint

.PHONY: server
server:
	VERSION=${VERSION} npm run dev

static/favicon.ico:
	convert static/apple-touch-icon.png -define icon:auto-resize=64,48,32,16 static/favicon.ico
