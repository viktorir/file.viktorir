CREATE TABLE IF NOT EXISTS files (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    size INTEGER NOT NULL,
    uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL,
    path TEXT NOT NULL,
    short_link TEXT NOT NULL UNIQUE,
    description TEXT,
    tags TEXT,
    hash TEXT NOT NULL,
    status TEXT DEFAULT 'active'
);