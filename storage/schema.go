package storage

const schema = `
BEGIN;

CREATE TABLE IF NOT EXISTS bookmarks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL,
	content TEXT NOT NULL,
	archived BOOLEAN NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS bookmarks_tags (
	bookmark_id INTEGER NOT NULL,
	name VARCHAR(64) NOT NULL,
	FOREIGN KEY(bookmark_id) REFERENCES bookmarks(id),
	UNIQUE(bookmark_id, name) ON CONFLICT IGNORE
);

CREATE VIRTUAL TABLE IF NOT EXISTS bookmarks_fts
USING fts5(title, url, content, content=bookmarks, content_rowid=id);

CREATE TRIGGER IF NOT EXISTS bookmarks_ai AFTER INSERT ON bookmarks BEGIN
	INSERT INTO bookmarks_fts(rowid, title, url, content) VALUES (new.id, new.title, new.url, new.content);
END;

CREATE TRIGGER IF NOT EXISTS bookmarks_ad AFTER DELETE ON bookmarks BEGIN
	INSERT INTO bookmarks_fts(bookmarks_fts, rowid, title, url, content) VALUES('delete', old.id, old.title, old.url, old.content);
END;

CREATE TRIGGER IF NOT EXISTS bookmarks_au AFTER UPDATE ON bookmarks BEGIN
	INSERT INTO bookmarks_fts(bookmarks_fts, rowid, title, url, content) VALUES('delete', old.id, old.title, old.url, old.content);
	INSERT INTO bookmarks_fts(rowid, title, url, content) VALUES (new.id, new.title, new.url, new.content);
END;

CREATE TABLE IF NOT EXISTS feeds (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	refreshed DATE DEFAULT (datetime('now')),
	last_authored DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	category VARCHAR(64) DEFAULT 'default',
	url VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS feeds_tags (
	feed_id INTEGER NOT NULL,
	name VARCHAR(64) NOT NULL,
	FOREIGN KEY(feed_id) REFERENCES feeds(id),
	UNIQUE(feed_id, name) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS items (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	feed_id INTEGER NOT NULL,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	date DATE DEFAULT (datetime('now')),
	url VARCHAR(255) UNIQUE NOT NULL,
	content TEXT NOT NULL,
	FOREIGN KEY(feed_id) REFERENCES feeds(id)
);

CREATE TABLE IF NOT EXISTS pages (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS pages_tags (
	page_id INTEGER NOT NULL,
	name VARCHAR(64) NOT NULL,
	FOREIGN KEY(page_id) REFERENCES pages(id),
	UNIQUE(page_id, name) ON CONFLICT IGNORE
);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(64) NOT NULL,
	password VARCHAR(255) NOT NULL,
	token VARCHAR(255) NOT NULL
);

CREATE VIEW IF NOT EXISTS tags AS
	SELECT name FROM bookmarks_tags
	UNION
	SELECT name FROM feeds_tags
	UNION
	SELECT name FROM pages_tags
;

COMMIT;
`

func (store *Store) migrate() error {
	_, err := store.db.Exec(schema)

	return err
}
