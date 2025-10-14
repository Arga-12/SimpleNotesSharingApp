-- notes table
CREATE TABLE IF NOT EXISTS notes (
  id SERIAL PRIMARY KEY,
  owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  content TEXT,
  shared BOOLEAN DEFAULT false,
  favorite BOOLEAN DEFAULT false,
  updated_at TIMESTAMP DEFAULT now()
);

-- index for owner
CREATE INDEX IF NOT EXISTS idx_notes_owner ON notes(owner_id);