CREATE TABLE IF NOT EXISTS summary(
  id uuid PRIMARY KEY NOT NULL,
  memory_id uuid NOT NULL,
  text TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  FOREIGN KEY(memory_id) REFERENCES memory(id) ON DELETE CASCADE
);
