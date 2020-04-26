package storage

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/timshannon/bolthold"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user that can authenticate to the bookmarks service
type User struct {
	Username string `boltholdKey:"Username"`
	Password string
	Token    string
	Created  time.Time
	Updated  time.Time
}

// GetUser finds a single user by URL
func (store *Store) GetUser(user *User) error {
	return store.db.FindOne(user, bolthold.Where(bolthold.Key).Eq(user.Username))
}

// Authenticate authenticates a user with a password and returns the token for the user
func (store *Store) Authenticate(username string, password string) (string, error) {
	user := User{
		Username: username,
	}

	if err := store.GetUser(&user); err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}

	return user.Token, nil
}

// UserTokenExists checks if there is exactly one user with the given token
func (store *Store) UserTokenExists(token string) bool {
	count, err := store.db.Count(&User{}, bolthold.Where("Token").Eq(token))
	if err != nil {
		log.Warn().Err(err).Msg("Error fetching token")
		return false
	}

	return count == 1
}
