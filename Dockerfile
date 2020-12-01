# syntax = docker/dockerfile:1-experimental
FROM golang:alpine AS gobase
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
        golang.org/x/lint/golint
WORKDIR /src



FROM node:alpine AS nodebase
RUN npm install -g @vue/cli
WORKDIR /src/web



FROM nodebase AS nodebuilder
COPY web/package*.json /src/web/
RUN yarn install --no-progress
COPY web/ /src/web/
RUN yarn run lint --no-progress
RUN yarn run build --no-progress --production



FROM gobase AS gobuilder
ENV CGO_ENABLED=1
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
COPY . ./
COPY --from=nodebuilder /src/web/dist/ ./web/dist/
RUN --mount=type=cache,target=/root/.cache/go-build go generate -v api/api.go
ARG VERSION=docker
ARG COMMIT=unknown
ARG DATE=unknown
RUN --mount=type=cache,target=/root/.cache/go-build go vet -v ./...
RUN --mount=type=cache,target=/root/.cache/go-build golint ./...
# RUN go test -v -covermode=count ./...
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o bookmarks \
        --tags "json1 fts5" \
        -ldflags "\
            -X github.com/nrocco/bookmarks/cmd.version=${VERSION} \
            -X github.com/nrocco/bookmarks/cmd.commit=${COMMIT} \
            -X github.com/nrocco/bookmarks/cmd.date=${DATE} \
            -s -w"



FROM alpine:edge
RUN apk add --no-cache \
        ca-certificates \
        sqlite \
    && true
COPY --from=gobuilder /src/bookmarks /usr/bin/bookmarks
EXPOSE 3000
WORKDIR /var/lib/bookmarks
VOLUME /var/lib/bookmarks
CMD ["bookmarks", “server”]
