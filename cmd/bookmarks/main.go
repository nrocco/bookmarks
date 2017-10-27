package main

import (
	"flag"
	"log"

	"github.com/nrocco/bookmarks/pkg/server"
)

var (
	// Workers stores the amount of workers that can do async tasks
	Workers = flag.Int("workers", 4, "The number of workers to start")

	Interval = flag.Int("interval", 30, "Fetch new feeds with this interval in minutes")

	// HTTPAddr stores the value for the --http option and defaults to 0.0.0.0:8000
	HTTPAddr = flag.String("http", "0.0.0.0:3000", "Address to listen for HTTP requests on")

	// Database holds the connection string for the database connection
	Database = flag.String("database", "data.sqlite", "The location to the sqlite database")
)

func main() {
	flag.Parse()

	log.Printf("%d workers are serving bookmarks from %s at http://%s\n", *Workers, *Database, *HTTPAddr)

	if err := server.Start(*Database, *HTTPAddr, *Workers, *Interval); err != nil {
		log.Fatal(err.Error())
	}
}
