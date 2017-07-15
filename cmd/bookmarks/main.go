package main

import (
	"flag"
	"github.com/nrocco/bookmarks/pkg/server"
	"log"
)

var (
	// HTTPAddr stores the value for the --http option and defaults to 0.0.0.0:8000
	HTTPAddr = flag.String("http", "0.0.0.0:8000", "Address to listen for HTTP requests on")

	// Database holds the connection string for the database connection
	Database = flag.String("database", "", "The connection string of the database server")

	// Secret is the value
	Secret = flag.String("secret", "", "The secret hash to authenticate to the api")
)

func main() {
	flag.Parse()

	app := server.App{
		Secret:           *Secret,
		ConnectionString: *Database,
	}

	if err := app.Run(*HTTPAddr); err != nil {
		log.Fatal(err.Error())
	}
}
