FROM golang:1.9-alpine AS builder
WORKDIR /go/src/github.com/nrocco/bookmarks
COPY . ./
RUN apk add --no-cache gcc musl-dev ca-certificates sqlite git
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure && dep status
RUN go generate pkg/server/app.go
RUN CGO_ENABLED=1 go build -x -v -a -o bookmarks github.com/nrocco/bookmarks/cmd/bookmarks

FROM alpine:edge
WORKDIR /var/lib/bookmarks
RUN apk add --no-cache ca-certificates sqlite
COPY --from=builder /go/src/github.com/nrocco/bookmarks/bookmarks /usr/bin/bookmarks
EXPOSE 3000
CMD ["bookmarks"]
