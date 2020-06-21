package storage

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nrocco/qb"

	// Store uses sqlite3 for its database
	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseFile = "data.db"
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.1 Safari/605.1.15"
)

func qbLogger(format string, v ...interface{}) {
}

// New returns a new instance of a Bookmarks Store
func New(path string) (*Store, error) {
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

func generateUUID() (uuid string) {
	b := make([]byte, 8)

	rand.Read(b)

	return strings.ToLower(fmt.Sprintf("%X", b))
}
