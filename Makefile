DOCKER_IMAGE = nrocco/bookmarks
DOCKER_IMAGE_VERSION = latest


.PHONY: container
container:
	docker image build \
		--build-arg "VERSION=$(shell git describe --tags --always --dirty)" \
		--build-arg "COMMIT=$(shell git describe --always --dirty)" \
		--build-arg "DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")" \
		--tag "$(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)" \
		.


.PHONY: push
push: container
	docker image push "$(DOCKER_IMAGE):$(DOCKER_IMAGE_VERSION)"
