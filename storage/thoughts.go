package storage

import (
	"regexp"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
)

// Thought holds information about a thought
type Thought struct {
	Created time.Time
	Updated time.Time
	Title   string `boltholdKey:"Title"`
	Content string
	Tags    []string `boltholdSliceIndex:"Tags"`
}

// ListThoughtsOptions can be passed to ListThoughts to filter thoughts
type ListThoughtsOptions struct {
	Search string
	Tags   []string
	Limit  int
	Offset int
}

// ListThoughts lists thoughts from the database
func (store *Store) ListThoughts(options *ListThoughtsOptions) (*[]*Thought, int) {
	thoughts := []*Thought{}

	query := &bolthold.Query{}

	if options.Search != "" {
		re, err := regexp.Compile("(?im)" + options.Search)
		if err != nil {
			return &thoughts, 0
		}

		query.And("Title").RegExp(re).Or(bolthold.Where("Content").RegExp(re))
	}

	for _, tag := range options.Tags {
		if tag == "" {
			continue
		} else if strings.HasPrefix(tag, "-") {
			query.And("Tags").Not().Contains(strings.TrimPrefix(tag, "-"))
		} else {
			query.And("Tags").Contains(tag)
		}
	}

	totalCount, err := store.db.Count(&Thought{}, query)
	if err != nil {
		log.Warn().Err(err).Msg("Error fetching thoughts count")
		return &thoughts, 0
	}

	query.Limit(options.Limit)
	query.Skip(options.Offset)

	if err := store.db.Find(&thoughts, query); err != nil {
		log.Warn().Err(err).Msg("Error fetching thoughts")
		return &thoughts, 0
	}

	return &thoughts, totalCount
}

// GetThought gets a single thought from the database
func (store *Store) GetThought(thought *Thought) error {
	return store.db.FindOne(thought, bolthold.Where(bolthold.Key).Eq(thought.Title))
}

// PersistThought adds a thought to the database
func (store *Store) PersistThought(thought *Thought) error {
	if thought.Created.IsZero() {
		thought.Created = time.Now()
	}

	thought.Updated = time.Now()

	return store.db.Upsert(thought.Title, thought)
}

// DeleteThought removes a thought from the database
func (store *Store) DeleteThought(thought *Thought) error {
	return store.db.Delete(thought.Title, thought)
}
