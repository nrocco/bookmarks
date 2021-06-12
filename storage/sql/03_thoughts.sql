CREATE TABLE IF NOT EXISTS thoughts (
    id CHAR(16) PRIMARY KEY,
    created DATE NOT NULL,
    updated DATE NOT NULL,
    tags JSON NOT NULL DEFAULT '[]',
    content TEXT NOT NULL DEFAULT ''
);

CREATE VIRTUAL TABLE IF NOT EXISTS thoughts_fts
USING fts5(content, tags, content=thoughts, content_rowid=rowid);

CREATE TRIGGER IF NOT EXISTS thoughts_ai AFTER INSERT ON thoughts BEGIN
    INSERT INTO thoughts_fts(rowid, content, tags) VALUES (new.rowid, new.content, new.tags);
END;

CREATE TRIGGER IF NOT EXISTS thoughts_ad AFTER DELETE ON thoughts BEGIN
    INSERT INTO thoughts_fts(thoughts_fts, rowid, content, tags) VALUES('delete', old.rowid, old.content, old.tags);
END;

CREATE TRIGGER IF NOT EXISTS thoughts_au AFTER UPDATE ON thoughts BEGIN
    INSERT INTO thoughts_fts(thoughts_fts, rowid, content, tags) VALUES('delete', old.rowid, old.content, old.tags);
    INSERT INTO thoughts_fts(rowid, content, tags) VALUES (new.rowid, new.content, new.tags);
END;
