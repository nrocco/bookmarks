package storage

// Tag represents a single tag
type Tag struct {
	Name  string
	Count int64
}

func (store *Store) ListTags() *[]*Tag {
	tags := []*Tag{}

	query := store.db.Select("tags")
	query.Columns("DISTINCT name")
	query.OrderBy("name", "ASC")
	query.Load(&tags)

	return &tags
}
