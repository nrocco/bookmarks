package storage

import (
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	// ErrNoThoughtTitle is returned if the Thought does not have a Title
	ErrNoThoughtTitle = errors.New("Missing Thought.Title")

	// ErrNoThoughtKey is returned if the Thought does not have a ID or Title
	ErrNoThoughtKey = errors.New("Missing Thought.ID or Thought.Title")
)

// Thought holds information about a thought
type Thought struct {
	ID      string
	Created time.Time
	Updated time.Time
	Title   string
	Content string `json:",omitempty"`
	Tags    Tags
}

// ListThoughtsOptions can be passed to ListThoughts to filter thoughts
type ListThoughtsOptions struct {
	Search string
	Tags   Tags
	Limit  int
	Offset int
}

// ListThoughts lists thoughts from the database
func (store *Store) ListThoughts(options *ListThoughtsOptions) (*[]*Thought, int) {
	query := store.db.Select("thoughts")

	if options.Search != "" {
		query.Where("title LIKE ? OR content LIKE ?", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	for _, tag := range options.Tags {
		if tag == "" {
			continue
		} else if strings.HasPrefix(tag, "-") {
			query.Where("NOT EXISTS (SELECT 1 FROM json_each(thoughts.tags) where json_each.value = ?)", strings.TrimPrefix(tag, "-"))
		} else {
			query.Where("EXISTS (SELECT 1 FROM json_each(thoughts.tags) where json_each.value = ?)", tag)
		}
	}

	thoughts := []*Thought{}
	totalCount := 0

	query.Columns("COUNT(id)")
	if err := query.LoadValue(&totalCount); err != nil {
		log.Warn().Err(err).Msg("Error fetching thought count")
		return &thoughts, 0
	}

	query.Columns("id", "created", "updated", "tags", "title")
	query.OrderBy("created", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	if _, err := query.Load(&thoughts); err != nil {
		log.Warn().Err(err).Msg("Error fetching thoughts")
		return &thoughts, 0
	}

	return &thoughts, totalCount
}

// GetThought gets a single thought from the database
func (store *Store) GetThought(thought *Thought) error {
	query := store.db.Select("thoughts")
	query.Limit(1)

	if thought.ID != "" {
		query.Where("id = ?", thought.ID)
	} else if thought.Title != "" {
		query.Where("title = ?", thought.Title)
	} else {
		return ErrNoThoughtKey
	}

	if err := query.LoadValue(&thought); err != nil {
		return err
	}

	return nil
}

// PersistThought adds a thought to the database
func (store *Store) PersistThought(thought *Thought) error {
	if thought.Title == "" {
		return ErrNoThoughtTitle
	}

	if thought.ID == "" {
		thought.ID = generateUUID()
	}

	if thought.Created.IsZero() {
		thought.Created = time.Now()
	}

	thought.Updated = time.Now()

	if thought.ID == "" {
		thought.ID = generateUUID()

		query := store.db.Insert("thoughts")
		query.Columns("id", "title", "created", "content", "tags", "updated")
		query.OnConflict("title", "content=excluded.content, tags=excluded.tags, updated=excluded.updated")
		// query.Returning("id") TODO this does not work in combination with on conflict
		query.Record(thought)

		if _, err := query.Exec(); err != nil {
			log.Error().Err(err).Str("id", thought.ID).Str("title", thought.Title).Msg("Error persisting thought")
			return err
		}
	} else {
		query := store.db.Update("thoughts")
		query.Set("content", thought.Content)
		query.Set("tags", thought.Tags)
		query.Set("title", thought.Title)
		query.Set("updated", thought.Updated)
		query.Where("id = ?", thought.ID)

		if _, err := query.Exec(); err != nil {
			log.Error().Err(err).Str("id", thought.ID).Str("title", thought.Title).Msg("Error updating thought")
			return err
		}
	}

	log.Info().Str("id", thought.ID).Str("title", thought.Title).Msg("Persisted thought")

	return nil
}

// DeleteThought removes a thought from the database
func (store *Store) DeleteThought(thought *Thought) error {
	if thought.ID == "" && thought.Title == "" {
		return ErrNoThoughtKey
	}

	query := store.db.Delete("thoughts")

	if thought.ID != "" {
		query.Where("id = ?", thought.ID)
	}

	if thought.Title != "" {
		query.Where("title = ?", thought.Title)
	}

	if _, err := query.Exec(); err != nil {
		log.Error().Err(err).Str("id", thought.ID).Str("title", thought.Title).Msg("Error deleting thought")
		return err
	}

	log.Info().Str("id", thought.ID).Str("title", thought.Title).Msg("Thought refreshed")

	return nil
}
