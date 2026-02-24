CREATE TABLE IF NOT EXISTS mock_endpoints (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    path TEXT NOT NULL,
    method TEXT NOT NULL,
    response_body TEXT,
    status_code INTEGER DEFAULT 200,
    content_type TEXT DEFAULT 'application/json',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
