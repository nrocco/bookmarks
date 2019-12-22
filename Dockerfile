FROM golang:alpine as gobase
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        git \
        musl-dev \
        sqlite
RUN env GO111MODULE=on go get -u github.com/cortesi/modd/cmd/modd
RUN go get -u github.com/kevinburke/go-bindata/...
RUN go get -u golang.org/x/lint/golint
WORKDIR /src



FROM node:alpine AS npmbase
WORKDIR /src



FROM npmbase AS npmbuilder
WORKDIR /src
COPY web/package*.json /src/
RUN npm install --no-progress
COPY web/ /src/
RUN npm run lint --no-progress
RUN npm run build --production --no-progress



FROM gobase AS gobuilder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN golint ./...
COPY . ./
COPY --from=npmbuilder /src/dist/ ./web/dist/
RUN go generate -v api/api.go
ARG VERSION=unknown
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
RUN go build -v -o bookmarks \
        --tags "fts5" \
        -ldflags "\
            -X github.com/nrocco/bookmarks/cmd.version=${VERSION} \
            -X github.com/nrocco/bookmarks/cmd.commit=${COMMIT} \
            -X github.com/nrocco/bookmarks/cmd.buildDate=${BUILD_DATE}"



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
