package storage

import (
	"github.com/nrocco/qb"

	// Store uses sqlite3 for its database
	_ "github.com/mattn/go-sqlite3"
)

func New(file string) (*Store, error) {
	var err error

	db, err := qb.Open(file)
	if err != nil {
		return &Store{}, err
	}

	store := Store{db}

	if err := store.migrate(); err != nil {
		return &Store{}, err
	}

	return &store, nil
}

type Store struct {
	db *qb.DB
}
