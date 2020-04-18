FROM golang:alpine as gobase
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        git \
        musl-dev \
        sqlite
RUN env GO111MODULE=on go get -u \
        github.com/cortesi/modd/cmd/modd \
        github.com/kevinburke/go-bindata/... \
        golang.org/x/lint/golint
WORKDIR /src



FROM node:alpine AS nodebase
RUN npm install -g @vue/cli
WORKDIR /src



FROM nodebase AS nodebuilder
WORKDIR /src
COPY web/package*.json /src/
RUN yarn install --no-progress
COPY web/ /src/
RUN yarn run lint --no-progress
RUN yarn run build --no-progress --production



FROM gobase AS gobuilder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=nodebuilder /src/dist/ ./web/dist/
RUN go generate -v api/api.go
ARG VERSION=unknown
ARG COMMIT=unknown
ARG DATE=unknown
RUN go vet ./...
RUN golint ./...
RUN go build -v -o bookmarks \
        --tags "fts5" \
        -ldflags "\
            -X github.com/nrocco/bookmarks/cmd.version=${VERSION} \
            -X github.com/nrocco/bookmarks/cmd.commit=${COMMIT} \
            -X github.com/nrocco/bookmarks/cmd.date=${DATE}"



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
