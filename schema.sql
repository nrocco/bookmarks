CREATE TABLE bookmarks (
    id SERIAL PRIMARY KEY,
    created timestamp DEFAULT current_timestamp,
    url text NOT NULL,
    name text,
    content text,
    UNIQUE(url)
);
