package storage

import (
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/timshannon/bolthold"
)

const (
	databaseFile = "data.db"
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.1 Safari/605.1.15"
)

// New returns a new instance of a Bookmarks Store
func New(path string) (*Store, error) {
	db, err := bolthold.Open(filepath.Join(path, databaseFile), 0664, nil)
	if err != nil {
		return &Store{}, err
	}

	store := Store{db, path}

	return &store, nil
}

// Store is used to persist Bookmark, Feed and Thougt
type Store struct {
	db *bolthold.Store
	fs string
}

// generateID generates a random ID of 8 character
func generateID() (uuid string) {
	b := make([]byte, 8)

	rand.Read(b)

	return strings.ToLower(fmt.Sprintf("%X", b))
}
