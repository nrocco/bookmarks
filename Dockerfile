# syntax = docker/dockerfile:1-experimental
FROM --platform=${BUILDPLATFORM} golang:alpine AS gobase
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        git \
        musl-dev \
        sqlite \
    && true
RUN env GO111MODULE=on go get -u \
        github.com/cortesi/modd/cmd/modd \
        github.com/kevinburke/go-bindata/... \
        golang.org/x/lint/golint \
        golang.org/x/tools/cmd/goimports \
    && true
WORKDIR /src



FROM --platform=${BUILDPLATFORM} node:alpine AS nodebase
RUN npm install -g @vue/cli
WORKDIR /src/web



FROM --platform=${BUILDPLATFORM} nodebase AS nodebuilder
COPY web/package*.json /src/web/
RUN yarn install --no-progress
COPY web/ /src/web/
RUN yarn run lint --no-progress
RUN yarn run build --no-progress --production



FROM --platform=${BUILDPLATFORM} gobase AS gobuilder
ENV CGO_ENABLED=0
COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
ARG BUILD_VERSION=master
ARG BUILD_COMMIT=unknown
ARG BUILD_DATE=now
ARG TARGETOS
ARG TARGETARCH
COPY . .
COPY --from=nodebuilder /src/web/dist/ ./web/dist/
RUN --mount=type=cache,target=/root/.cache/go-build go generate -v api/api.go
RUN --mount=type=cache,target=/root/.cache/go-build golint -set_exit_status ./...
RUN --mount=type=cache,target=/root/.cache/go-build go vet -v ./...
RUN mkdir -p dist
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -x -o dist \
        -ldflags "\
            -X github.com/nrocco/bookmarks/cmd.version=${BUILD_VERSION} \
            -X github.com/nrocco/bookmarks/cmd.commit=${BUILD_COMMIT} \
            -X github.com/nrocco/bookmarks/cmd.date=${BUILD_DATE} \
            -s -w"
RUN --mount=type=cache,target=/root/.cache/go-build go test -v -short ./...



FROM scratch AS bin
COPY --from=gobuilder /src/dist/ /



FROM alpine:edge
RUN apk add --no-cache \
        ca-certificates \
        sqlite \
    && true
COPY --from=gobuilder /src/dist/bookmarks /usr/bin/bookmarks
EXPOSE 3000
WORKDIR /var/lib/bookmarks
VOLUME /var/lib/bookmarks
CMD ["bookmarks", "server"]
