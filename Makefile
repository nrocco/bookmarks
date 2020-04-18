DOCKER_IMAGE = nrocco/bookmarks

lint:
	golint ./...
	go vet ./...

container:
	docker build \
		--build-arg "VERSION=$(shell git describe --tags)" \
		--build-arg "COMMIT=$(shell git describe --always)" \
		--build-arg "DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)" \
		--tag "$(DOCKER_IMAGE):$(shell git describe --tags)" \
		--tag "$(DOCKER_IMAGE):latest" \
		.

push:
	docker push "$(DOCKER_IMAGE):$(shell git describe --tags)"
	docker push "$(DOCKER_IMAGE):latest"
