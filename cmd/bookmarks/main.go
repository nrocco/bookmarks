package main

import (
	"flag"
	"log"

	"github.com/nrocco/bookmarks/pkg/server"
)

var (
	// HTTPAddr stores the value for the --http option and defaults to 0.0.0.0:8000
	HTTPAddr = flag.String("http", "0.0.0.0:3000", "Address to listen for HTTP requests on")

	// Database holds the connection string for the database connection
	Database = flag.String("database", "data.sqlite", "The location to the sqlite database")
)

func main() {
	flag.Parse()

	log.Printf("Serving bookmarks from %s at http://%s\n", *Database, *HTTPAddr)

	if err := server.Start(*Database, *HTTPAddr); err != nil {
		log.Fatal(err.Error())
	}
}
