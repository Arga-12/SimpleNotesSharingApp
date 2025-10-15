-- logs table for request/response logging
CREATE TABLE IF NOT EXISTS logs (
  id SERIAL PRIMARY KEY,
  datetime TIMESTAMP DEFAULT NOW(),
  method VARCHAR(10) NOT NULL,
  endpoint TEXT NOT NULL,
  request_headers JSONB,
  request_payload TEXT,
  response_body TEXT,
  response_status INTEGER NOT NULL,
  duration_ms INTEGER,
  user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

-- Index untuk query performa
CREATE INDEX IF NOT EXISTS idx_logs_datetime ON logs(datetime DESC);
CREATE INDEX IF NOT EXISTS idx_logs_endpoint ON logs(endpoint);
CREATE INDEX IF NOT EXISTS idx_logs_status ON logs(response_status);
CREATE INDEX IF NOT EXISTS idx_logs_user ON logs(user_id);