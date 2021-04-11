package storage

import (
	"context"
	"crypto/rand"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nrocco/qb"

	// Store uses sqlite for its database
	_ "modernc.org/sqlite"
)

const (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.1 Safari/605.1.15"
)

// New returns a new instance of a Bookmarks Store
func New(ctx context.Context, path string) (*Store, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return &Store{}, err
	}

	db, err := qb.Open(ctx, path)
	if err != nil {
		return &Store{}, err
	}

	store := Store{db}

	if err := store.migrate(ctx); err != nil {
		return &Store{}, err
	}

	return &store, nil
}

// Store is used to persist Bookmark, Feed and Thought's
type Store struct {
	db *qb.DB
}

func generateUUID() (uuid string) {
	b := make([]byte, 8)

	rand.Read(b)

	return strings.ToLower(fmt.Sprintf("%X", b))
}
