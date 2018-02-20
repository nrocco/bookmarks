FROM node:alpine as npmbuilder
ARG VERSION
WORKDIR /src
COPY frontend ./
RUN npm install --no-progress
RUN npm run lint --no-progress
RUN npm run build --production --no-progress

FROM golang:alpine AS gobuilder
ARG VERSION
WORKDIR /go/src/github.com/nrocco/bookmarks
COPY . ./
COPY --from=npmbuilder /src/dist ./frontend/dist/
RUN apk add --no-cache gcc musl-dev ca-certificates sqlite git
RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/jteeuwen/go-bindata/...
RUN dep ensure && dep status
RUN go generate github.com/nrocco/bookmarks/...
RUN go build -v -o build/bookmarks -ldflags "-X main.Version=${VERSION}" github.com/nrocco/bookmarks/cmd/bookmarks

FROM alpine:edge
WORKDIR /var/lib/bookmarks
RUN apk add --no-cache ca-certificates sqlite
COPY --from=gobuilder /go/src/github.com/nrocco/bookmarks/build/bookmarks /usr/bin/bookmarks
EXPOSE 3000
CMD ["bookmarks"]
