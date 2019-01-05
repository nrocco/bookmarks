FROM node:alpine AS npmbuilder
RUN mkdir -p /app
WORKDIR /app
COPY web/package*.json /app/
RUN npm install --no-progress
COPY web/ /app/
RUN npm run lint --no-progress
RUN npm run build --production --no-progress
###############################################################################
FROM golang:alpine AS gobuilder
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        git \
        musl-dev \
        sqlite \
    && go get -u github.com/jteeuwen/go-bindata/... \
    && go get -u github.com/cortesi/modd/cmd/modd \
    && go get -u golang.org/x/lint/golint
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN golint ./...
COPY . ./
COPY --from=npmbuilder /app/dist/ ./web/dist/
RUN go generate -v api/api.go && \
    go generate -v storage/bookmarks.go
ARG VERSION=unknown
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
RUN go build -v -o bookmarks \
        -ldflags "\
            -X github.com/nrocco/bookmarks/cmd.version=${VERSION} \
            -X github.com/nrocco/bookmarks/cmd.commit=${COMMIT} \
            -X github.com/nrocco/bookmarks/cmd.buildDate=${BUILD_DATE}"
###############################################################################
FROM alpine:edge
RUN apk add --no-cache \
        ca-certificates \
        sqlite \
    && true
COPY --from=gobuilder /src/bookmarks /usr/bin/bookmarks
EXPOSE 3000
WORKDIR /var/lib/bookmarks
VOLUME /var/lib/bookmarks
CMD ["bookmarks"]
