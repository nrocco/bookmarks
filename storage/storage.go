package storage

import (
	"path/filepath"

	"github.com/nrocco/qb"

	// Store uses sqlite3 for its database
	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseFile = "data.sqlite"
)

func New(path string) (*Store, error) {
	var err error

	path, err = filepath.Abs(path)
	if err != nil {
		return &Store{}, err
	}

	db, err := qb.Open(filepath.Join(path, databaseFile))
	if err != nil {
		return &Store{}, err
	}

	store := Store{db, path}

	if err := store.migrate(); err != nil {
		return &Store{}, err
	}

	return &store, nil
}

type Store struct {
	db *qb.DB
	fs string
}
