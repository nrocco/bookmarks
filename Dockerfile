FROM node:latest AS npmbuilder
RUN mkdir -p /app
WORKDIR /app
COPY web/package*.json /app/
RUN npm install --no-progress
COPY web/ /app/
RUN npm run lint --no-progress
RUN npm run build --production --no-progress

FROM golang:alpine AS gobuilder
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        git \
        musl-dev \
        sqlite \
    && go get -u github.com/jteeuwen/go-bindata/... \
    && go get github.com/cortesi/modd/cmd/modd
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
COPY --from=npmbuilder /app/dist/ ./web/dist/
RUN go generate -v api/api.go && \
    go generate -v storage/bookmarks.go
ARG VERSION=unknown
RUN go build -v -o bookmarks \
        -ldflags "-X main.Version=${VERSION}" \
        cmd/bookmarks/bookmarks.go

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
