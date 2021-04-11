package storage

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Thought holds information about a thought
type Thought struct {
	ID      string
	Created time.Time
	Updated time.Time
	Content string
	Tags    Tags
}

// ThoughtListOptions can be passed to ThoughtList to filter thoughts
type ThoughtListOptions struct {
	Search string
	Tags   Tags
	Limit  int
	Offset int
}

// ThoughtList lists thoughts from the database
func (store *Store) ThoughtList(ctx context.Context, options *ThoughtListOptions) (*[]*Thought, int) {
	query := store.db.Select(ctx).From("thoughts")

	if options.Search != "" {
		query.Where("rowid IN (SELECT rowid FROM thoughts_fts(?))", options.Search)
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
		log.Ctx(ctx).Warn().Err(err).Msg("Error fetching thought count")
		return &thoughts, 0
	}

	query.Columns("id", "created", "updated", "content", "tags")
	query.OrderBy("created", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	if _, err := query.Load(&thoughts); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Error fetching thoughts")
		return &thoughts, 0
	}

	return &thoughts, totalCount
}

// ThoughtGet gets a single thought from the database
func (store *Store) ThoughtGet(ctx context.Context, thought *Thought) error {
	query := store.db.Select(ctx).From("thoughts")
	query.Limit(1)

	if thought.ID != "" {
		query.Where("id = ?", thought.ID)
	} else {
		return errors.New("Missing Thought.ID")
	}

	if err := query.LoadValue(&thought); err != nil {
		return err
	}

	return nil
}

// ThoughtPersist adds a thought to the database
func (store *Store) ThoughtPersist(ctx context.Context, thought *Thought) error {
	if thought.Created.IsZero() {
		thought.Created = time.Now()
	}

	if thought.Tags == nil {
		thought.Tags = Tags{}
	}

	thought.Updated = time.Now()

	if thought.ID == "" {
		thought.ID = generateUUID()

		query := store.db.Insert(ctx).InTo("thoughts")
		query.Columns("id", "created", "content", "tags", "updated")
		query.Record(thought)

		if _, err := query.Exec(); err != nil {
			log.Ctx(ctx).Error().Err(err).Str("id", thought.ID).Msg("Error persisting thought")
			return err
		}
	} else {
		query := store.db.Update(ctx).Table("thoughts")
		query.Set("content", thought.Content)
		query.Set("tags", thought.Tags)
		query.Set("updated", thought.Updated)
		query.Where("id = ?", thought.ID)

		if _, err := query.Exec(); err != nil {
			log.Ctx(ctx).Error().Err(err).Str("id", thought.ID).Msg("Error updating thought")
			return err
		}
	}

	log.Ctx(ctx).Info().Str("id", thought.ID).Msg("Persisted thought")

	return nil
}

// ThoughtDelete removes a thought from the database
func (store *Store) ThoughtDelete(ctx context.Context, thought *Thought) error {
	if thought.ID == "" {
		return errors.New("Missing Thought.ID")
	}

	query := store.db.Delete(ctx).From("thoughts")

	if thought.ID != "" {
		query.Where("id = ?", thought.ID)
	}

	if _, err := query.Exec(); err != nil {
		log.Ctx(ctx).Error().Err(err).Str("id", thought.ID).Msg("Error deleting thought")
		return err
	}

	log.Ctx(ctx).Info().Str("id", thought.ID).Msg("Thought refreshed")

	return nil
}
