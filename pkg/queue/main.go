package queue

import (
	log "github.com/sirupsen/logrus"

	"github.com/nrocco/bookmarks/pkg/storage"
)

// New returns a buffered channel that we can send work requests on.
func New(store *storage.Store, nworkers int) *Queue {
	queue := Queue{
		work:    make(chan workRequest, 100),
		workers: make(chan chan workRequest, nworkers),
	}

	for i := 0; i < nworkers; i++ {
		worker := worker{
			store:   store,
			work:    make(chan workRequest),
			workers: queue.workers,
		}
		worker.Start()
	}

	queue.start()

	return &queue
}

type Queue struct {
	work    chan workRequest
	workers chan chan workRequest
}

func (q *Queue) start() {
	go func() {
		for {
			select {
			case work := <-q.work:
				go func() {
					worker := <-q.workers
					worker <- work
				}()
			}
		}
	}()
}

func (q *Queue) Schedule(workType string, ID int64) {
	q.work <- workRequest{Type: workType, ID: ID}
}

type worker struct {
	store   *storage.Store
	work    chan workRequest
	workers chan chan workRequest
}

func (w *worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.workers <- w.work

			select {
			case work := <-w.work:
				l := log.WithFields(log.Fields{
					"type": work.Type,
					"id":   work.ID,
				})

				if work.Type == "Bookmark.FetchContent" {
					bookmark := storage.Bookmark{ID: work.ID}
					if err := w.store.GetBookmark(&bookmark); err != nil {
						l.WithError(err).Warn("Error loading bookmark")
						return
					}

					if err := bookmark.FetchContent(); err != nil {
						l.WithError(err).Warn("Error fetching content")
						return
					}

					if err := w.store.UpdateBookmark(&bookmark); err != nil {
						l.WithError(err).Warn("Error saving content")
						return
					}
				} else if work.Type == "Feed.Refresh" {
					feed := storage.Feed{ID: work.ID}
					if err := w.store.GetFeed(&feed); err != nil {
						l.WithError(err).Warn("Error loading feed")
						return
					}

					if err := w.store.RefreshFeed(&feed); err != nil {
						l.WithError(err).Warn("Error refreshing feed")
						return
					}
				} else {
					log.Printf("Unknown work received: %s", work.Type)
				}
			}
		}
	}()
}

type workRequest struct {
	Type string
	ID   int64
}
