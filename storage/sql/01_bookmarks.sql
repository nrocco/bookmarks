CREATE TABLE IF NOT EXISTS bookmarks (
    id CHAR(16) PRIMARY KEY,
    created DATE DEFAULT (datetime('now')),
    updated DATE DEFAULT (datetime('now')),
    title VARCHAR(64) NOT NULL,
    url VARCHAR(255) UNIQUE NOT NULL,
    excerpt TEXT NOT NULL DEFAULT '',
    content TEXT NOT NULL DEFAULT '',
    tags JSON NOT NULL DEFAULT '[]'
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
    INSERT INTO bookmarks_fts(rowid, title, url, content) VALUES (new.rowid, new.title, new.url, new.content);
END;
