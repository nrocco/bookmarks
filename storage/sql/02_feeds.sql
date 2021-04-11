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
