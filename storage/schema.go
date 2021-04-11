package storage

import (
	"context"
	"embed"
)

//go:embed sql/*.sql
var migrations embed.FS

func (store *Store) migrate(ctx context.Context) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	files, err := migrations.ReadDir("sql")
	if err != nil {
		return err
	}

	for _, file := range files {
		migration, err := migrations.ReadFile("sql/" + file.Name())
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, string(migration[:])); err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return err
}
