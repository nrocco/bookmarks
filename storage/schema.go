package storage

const schema = `
BEGIN;

CREATE TABLE IF NOT EXISTS bookmarks (
	id CHAR(16) PRIMARY KEY,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL,
	excerpt TEXT NOT NULL DEFAULT '',
	content TEXT NOT NULL DEFAULT '',
	tags JSON NOT NULL DEFAULT '[]',
	archived BOOLEAN NOT NULL DEFAULT 0
);

CREATE VIRTUAL TABLE IF NOT EXISTS bookmarks_fts
USING fts5(title, url, content, content=bookmarks, content_rowid=rowid);

CREATE TRIGGER IF NOT EXISTS bookmarks_ai AFTER INSERT ON bookmarks BEGIN
	INSERT INTO bookmarks_fts(rowid, title, url, content) VALUES (new.rowid, new.title, new.url, new.content);
END;

CREATE TRIGGER IF NOT EXISTS bookmarks_ad AFTER DELETE ON bookmarks BEGIN
	INSERT INTO bookmarks_fts(bookmarks_fts, rowid, title, url, content) VALUES('delete', old.rowid, old.title, old.url, old.content);
END;

CREATE TRIGGER IF NOT EXISTS bookmarks_au AFTER UPDATE ON bookmarks BEGIN
	INSERT INTO bookmarks_fts(bookmarks_fts, rowid, title, url, content) VALUES('delete', old.rowid, old.title, old.url, old.content);
	INSERT INTO bookmarks_fts(rowid, title, url, content) VALUES (new.id, new.title, new.url, new.content);
END;

CREATE TABLE IF NOT EXISTS feeds (
	id CHAR(16) PRIMARY KEY,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	refreshed DATE DEFAULT (datetime('now')),
	last_authored DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL,
	etag VARCHAR(200) NOT NULL DEFAULT '',
	tags JSON NOT NULL DEFAULT '[]',
	items JSON NOT NULL DEFAULT '[]'
);

CREATE TABLE IF NOT EXISTS users (
	id CHAR(16) PRIMARY KEY,
	username VARCHAR(64) NOT NULL,
	password VARCHAR(255) NOT NULL,
	token VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS thoughts (
	id CHAR(16) PRIMARY KEY,
	created DATE NOT NULL,
	updated DATE NOT NULL,
	title VARCHAR(255) UNIQUE NOT NULL,
	tags JSON NOT NULL DEFAULT '[]',
	content TEXT NOT NULL DEFAULT ''
);

COMMIT;
`

func (store *Store) migrate() error {
	_, err := store.db.Exec(schema)

	return err
}
