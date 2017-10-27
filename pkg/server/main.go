package server

import (
	"net/http"
)

//go:generate go-bindata -pkg server -o bindata.go assets

// Start initializes and starts all services
func Start(file string, address string, workers int, interval int) error {
	if err := initDB(file); err != nil {
		return err
	}

	initQueue(workers)

	initScheduler(interval)

	return http.ListenAndServe(address, initRouter())
}
