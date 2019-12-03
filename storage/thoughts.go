package storage

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Thought holds information about a thought
type Thought struct {
	ID      int64
	Created time.Time
	Updated time.Time
	Title   string
	Type    string
	Size    int64
	Content string
	Tags    []string
}

// Path returns the path to the thought on disk
func (thought *Thought) Path() string {
	return filepath.Join("var", "thoughts", thought.Title)
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

	query.Columns("t.id", "t.created", "t.updated", "t.title", "t.type", "t.size", "t.content")
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
func (store *Store) PersistThought(thought *Thought, data io.ReadCloser) error {
	if data == nil && thought.ID == 0 {
		return errors.New("Fuuuuu")
	}

	if data != nil {
		outFile, err := os.Create(thought.Path())
		if err != nil {
			return err
		}

		defer outFile.Close()

		thought.Size, err = io.Copy(outFile, data)
		if err != nil {
			return err
		}

		outFile.Seek(0, 0)

		content, err := ioutil.ReadAll(outFile)
		if err != nil {
			return err
		}

		thought.Content = string(content)
		thought.Type = "text/plain" // TODO: fix this

		outFile.Close()
	}

	thought.Updated = time.Now()

	if thought.ID == 0 {
		thought.Created = time.Now()

		query := store.db.Insert("thoughts")
		query.Columns("created", "updated", "title", "type", "size", "content")
		query.Record(thought)

		if _, err := query.Exec(); err != nil {
			return err
		}
	} else {
		query := store.db.Update("thoughts")
		query.Set("updated", thought.Updated)

		// query.Set("title", thought.Title) TODO: support changing titles

		if data != nil {
			query.Set("type", thought.Type)
			query.Set("size", thought.Size)
			query.Set("content", thought.Content)
		}

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

	if err := os.Remove(thought.Path()); err != nil {
		return err
	}

	query := store.db.Delete("thoughts")
	query.Where("id = ?", thought.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}
