package queue

import (
	"github.com/nrocco/bookmarks/storage"
	"github.com/rs/zerolog/log"
)

// New returns a buffered channel that we can send work requests on.
func New(store *storage.Store, nworkers int) *Queue {
	log.Info().Msg("Setting up the queue")

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
		log.Info().Int("worker", i).Msg("Starting worker")
		worker.Start()
	}

	queue.start()

	log.Info().Msg("Queue started")

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
	log.Info().Int64("id", ID).Str("work_type", workType).Msg("Scheduling work")

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
				logger := log.With().Int64("id", work.ID).Str("type", work.Type).Logger()

				if work.Type == "Bookmark.FetchContent" {
					bookmark := storage.Bookmark{ID: work.ID}
					if err := w.store.GetBookmark(&bookmark); err != nil {
						logger.Warn().Err(err).Msg("Error loading bookmark")
						return
					}

					if err := bookmark.FetchContent(); err != nil {
						logger.Warn().Err(err).Msg("Error fetching content")
						return
					}

					if err := w.store.UpdateBookmark(&bookmark); err != nil {
						logger.Warn().Err(err).Msg("Error saving content")
						return
					}
				} else if work.Type == "Feed.Refresh" {
					feed := storage.Feed{ID: work.ID}
					if err := w.store.GetFeed(&feed); err != nil {
						logger.Warn().Err(err).Msg("Error loading feed")
						return
					}

					if err := w.store.RefreshFeed(&feed); err != nil {
						logger.Warn().Err(err).Msg("Error refreshing feed")
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
