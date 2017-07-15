FROM golang:1.9 AS builder
WORKDIR /go/src/github.com/nrocco/bookmarks
COPY Gopkg.lock Gopkg.toml ./
COPY pkg ./pkg
COPY cmd ./cmd
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
ENV CGO_ENABLED=0
RUN go build -a -tags netgo -installsuffix netgo -ldflags "-d -s -w" -o build/bookmarks github.com/nrocco/bookmarks/cmd/bookmarks

FROM scratch
WORKDIR /app
COPY --from=builder /go/src/github.com/nrocco/bookmarks/build/bookmarks /app
COPY public /app/public/
ADD ca-certificates.crt /etc/ssl/certs/
EXPOSE 8000
CMD ["/app/bookmarks"]
