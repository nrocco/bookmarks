lint:
	golint ./...
	go vet ./...

container:
	docker build \
		--build-arg "VERSION=$(shell git describe --tags)" \
		--build-arg "COMMIT=$(shell git describe --always)" \
		--build-arg "BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" \
		--tag "nrocco/bookmarks:$(shell git describe --tags)" \
		--tag "nrocco/bookmarks:latest" \
		.

push:
	docker push "nrocco/bookmarks:$(shell git describe --tags)"
	docker push "nrocco/bookmarks:latest"
