package server

const schema = `CREATE TABLE IF NOT EXISTS bookmarks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	created DATE DEFAULT (datetime('now')),
	updated DATE DEFAULT (datetime('now')),
	title VARCHAR(64) NOT NULL,
	url VARCHAR(255) UNIQUE NOT NULL,
	content TEXT NOT NULL,
	archived BOOLEAN NOT NULL DEFAULT 0
);`
