package storage

func (store *Store) UserTokenExists(token string) bool {
	var count int64

	query := store.db.Select("users")
	query.Columns("COUNT(id)")
	query.Where("token = ?", token)
	query.LoadValue(&count)

	return count == 1
}

func (store *Store) UserPasswordHash(username string) string {
	var password string

	query := store.db.Select("users")
	query.Columns("password")
	query.Where("username = ?", username)
	query.LoadValue(&password)

	return password
}

func (store *Store) UserToken(username string) string {
	var token string

	query := store.db.Select("users")
	query.Columns("token")
	query.Where("username = ?", username)
	query.LoadValue(&token)

	return token
}
