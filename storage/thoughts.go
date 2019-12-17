package storage

import (
	"errors"
	"strings"
	"time"
)

// Thought holds information about a thought
type Thought struct {
	ID      int64 `json:"-"`
	Created time.Time
	Updated time.Time
	Title   string
	Content string
	Tags    []string
}

// ListThoughtsOptions can be passed to ListThoughts to filter thoughts
type ListThoughtsOptions struct {
	Search string
	Tags   string
	Limit  int
	Offset int
}

// ListThoughts lists thoughts from the database
func (store *Store) ListThoughts(options *ListThoughtsOptions) (*[]*Thought, int) {
	query := store.db.Select("thoughts t")

	if options.Search != "" {
		query.Where("t.title LIKE ? OR t.content LIKE ?", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	for _, tag := range strings.Split(options.Tags, ",") {
		if tag == "" {
			continue
		}

		if strings.HasPrefix(tag, "-") {
			query.Where("t.id NOT IN (SELECT id FROM thoughts_tags WHERE name = ?)", strings.TrimPrefix(tag, "-"))
		} else {
			query.Where("t.id IN (SELECT id FROM thoughts_tags WHERE name = ?)", tag)
		}
	}

	thoughts := []*Thought{}
	totalCount := 0

	query.Columns("COUNT(t.id)")
	query.LoadValue(&totalCount)

	query.Columns("t.id", "t.created", "t.updated", "t.title", "t.content")
	query.OrderBy("t.created", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	query.Load(&thoughts)

	for _, thought := range thoughts {
		thought.Tags = []string{}

		query = store.db.Select("thoughts_tags")
		query.Columns("name")
		query.Where("id = ?", thought.ID)
		query.LoadValue(&thought.Tags)
	}

	return &thoughts, totalCount
}

// GetThought gets a single thought from the database
func (store *Store) GetThought(thought *Thought) error {
	query := store.db.Select("thoughts")
	query.Limit(1)

	query.Where("title = ?", thought.Title)

	if err := query.LoadValue(&thought); err != nil {
		return err
	}

	thought.Tags = []string{}

	query = store.db.Select("thoughts_tags")
	query.Columns("name")
	query.Where("id = ?", thought.ID)
	query.LoadValue(&thought.Tags)

	return nil
}

// PersistThought adds a thought to the database
func (store *Store) PersistThought(thought *Thought) error {
	if thought.ID == 0 {
		thought.Created = time.Now()
		thought.Updated = time.Now()

		query := store.db.Insert("thoughts")
		query.Columns("created", "updated", "title", "content")
		query.Record(thought)

		if _, err := query.Exec(); err != nil {
			return err
		}
	} else {
		thought.Updated = time.Now()

		query := store.db.Update("thoughts")
		query.Set("updated", thought.Updated)
		query.Set("title", thought.Title)
		query.Set("content", thought.Content)
		query.Where("id = ?", thought.ID)

		if _, err := query.Exec(); err != nil {
			return err
		}
	}

	for _, tag := range thought.Tags {
		query := store.db.Insert("thoughts_tags").OrIgnore()
		query.Columns("id", "name")
		query.Values(thought.ID, tag)

		if _, err := query.Exec(); err != nil {
			return err
		}
	}

	query := store.db.Delete("thoughts_tags")
	query.Where("id = ?", thought.ID)

	for _, tag := range thought.Tags {
		query.Where("name != ?", tag)
	}

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}

// DeleteThought removes a thought from the database
func (store *Store) DeleteThought(thought *Thought) error {
	if thought.ID == 0 {
		return errors.New("Not an existing thought")
	}

	query := store.db.Delete("thoughts")
	query.Where("id = ?", thought.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}
