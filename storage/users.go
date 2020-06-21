package storage

import (
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user that can authenticate to the bookmarks service
type User struct {
	Username string
	Password string
	Token    string
	Created  time.Time
	Updated  time.Time
}

// GetUser finds a single user by Username
func (store *Store) GetUser(user *User) error {
	return store.db.Select("users").Columns("*").Where("username = ?", user.Username).Limit(1).LoadValue(user)
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
	count := 0

	query := store.db.Select("users").Columns("COUNT(*)").Where("token = ?", token)
	if err := query.LoadValue(&count); err != nil {
		log.Warn().Err(err).Msg("Error fetching token")
		return false
	}

	return count == 1
}
