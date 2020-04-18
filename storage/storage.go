package storage

import (
	"fmt"
	"path/filepath"

	"github.com/nrocco/qb"

	// Store uses sqlite3 for its database
	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseFile = "data.sqlite"
)

func qbLogger(format string, v ...interface{}) {
}

// New returns a new instace of a Bookmarks Store
func New(path string) (*Store, error) {
	var err error

	path, err = filepath.Abs(path)
	if err != nil {
		return &Store{}, err
	}

	db, err := qb.Open(filepath.Join(path, fmt.Sprintf("%s?_foreign_keys=yes", databaseFile)), qbLogger)
	if err != nil {
		return &Store{}, err
	}

	store := Store{db, path}

	if err := store.migrate(); err != nil {
		return &Store{}, err
	}

	return &store, nil
}

// Store is used to persist Bookmark, Feed and Thougt
type Store struct {
	db *qb.DB
	fs string
}
