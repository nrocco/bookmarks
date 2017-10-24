package server

const schema = `
CREATE TABLE IF NOT EXISTS bookmarks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL,
	content TEXT NOT NULL,
	archived BOOLEAN NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS feeds (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	refreshed DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	subtitle VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL
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
`
