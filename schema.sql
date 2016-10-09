CREATE TABLE bookmarks (
    id SERIAL PRIMARY KEY,
    created timestamp DEFAULT current_timestamp,
    url text NOT NULL,
    name text,
    content text,
    fts tsvector,
    UNIQUE(url)
);
CREATE INDEX bookmarks_fts_idx ON bookmarks USING gin(fts);
