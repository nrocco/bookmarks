package main

import (
	"flag"
	"github.com/nrocco/bookmarks/pkg/server"
	"log"
)

var (
	// HTTPAddr stores the value for the --http option and defaults to 0.0.0.0:8000
	HTTPAddr = flag.String("http", "0.0.0.0:8000", "Address to listen for HTTP requests on")

	// Host stores the value for the --db-host option
	Host = flag.String("db-host", "localhost", "The host of the database server to connect to")

	// User stores the value for the --db-user option
	User = flag.String("db-user", "postgres", "The username of the database server to connect to")

	// Pass stores the value for the --db-pass option
	Pass = flag.String("db-pass", "secret", "The password of the database server to connect to")

	// Name stores the value for the --db-name option
	Name = flag.String("db-name", "bookmarks", "The name of the database to connect to")

	Secret = flag.String("secret", "secret", "The secret hash to authenticate to the api")
)

func main() {
	flag.Parse()

	app := server.App{
		Secret: *Secret,
	}

	if err := app.Initialize(*Host, *User, *Pass, *Name); err != nil {
		log.Fatal(err.Error())
	}

	if err := app.Run(*HTTPAddr); err != nil {
		log.Fatal(err.Error())
	}
}
