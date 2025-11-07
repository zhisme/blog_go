-- Create mailing_list table
CREATE TABLE IF NOT EXISTS mailing_list (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE UNIQUE INDEX IF NOT EXISTS idx_mailing_list_email ON mailing_list(email);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_mailing_list_created_at ON mailing_list(created_at);
