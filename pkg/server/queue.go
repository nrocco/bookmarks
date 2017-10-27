package server

import (
	"log"
)

type WorkRequest struct {
	Type     string
	Bookmark Bookmark
	Feed     Feed
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
}

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 100)

var WorkerQueue chan chan WorkRequest

func initQueue(nworkers int) {
	log.Printf("Starting the queue with %d workers", nworkers)

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)

	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
		// log.Printf("Worker %d is ready", worker.ID)
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()
}

func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	return Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				if work.Type == "Bookmark.FetchContent" {
					if err := work.Bookmark.FetchContent(); err != nil {
						log.Println(err)
					}
				} else if work.Type == "Feed.Refresh" {
					if err := work.Feed.Refresh(); err != nil {
						log.Println(err)
					}
				} else {
					log.Printf("Unknown type of work received: %s", work.Type)
				}
			}
		}
	}()
}
