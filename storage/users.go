package storage

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user that can authenticate to the bookmarks service
type User struct {
	ID       string    `json:"id"`
	Created  time.Time `json:"-"`
	Updated  time.Time `json:"-"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Token    string    `json:"token"`
}

func (a *User) SetPassword(password string) error {
	var encryptedPassword []byte
	var err error

	encryptedPassword, err = bcrypt.GenerateFromPassword([]byte(password), 1)
	if err != nil {
		return err
	}

	a.Password = string(encryptedPassword[:])

	return nil
}

func (store *Store) UserAdd(ctx context.Context, user *User) error {
	user.ID = generateUUID()
	user.Created = time.Now()
	user.Updated = time.Now()

	query := store.db.Insert(ctx).InTo("users")
	query.Columns("id", "created", "updated", "username", "password", "token")
	query.Record(user)

	if _, err := query.Exec(); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Could not create user")
		return err
	}

	log.Ctx(ctx).Info().Str("id", user.ID).Str("username", user.Username).Msg("User created")

	return nil
}

func (store *Store) UserTokenExists(ctx context.Context, token string) bool {
	var count int64

	query := store.db.Select(ctx).From("users")
	query.Columns("COUNT(id)")
	query.Where("token = ?", token)
	query.LoadValue(&count)

	return count == 1
}

func (store *Store) UserPasswordHash(ctx context.Context, username string) string {
	var password string

	query := store.db.Select(ctx).From("users")
	query.Columns("password")
	query.Where("username = ?", username)
	query.LoadValue(&password)

	return password
}

func (store *Store) UserTokenGet(ctx context.Context, username string) string {
	var token string

	query := store.db.Select(ctx).From("users")
	query.Columns("token")
	query.Where("username = ?", username)
	query.LoadValue(&token)

	return token
}
