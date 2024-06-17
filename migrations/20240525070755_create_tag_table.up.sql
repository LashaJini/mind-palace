CREATE TABLE IF NOT EXISTS {{ .Namespace }}.tag(
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS {{ .Namespace }}.memory_tag(
  tag_id INTEGER NOT NULL,
  memory_id uuid NOT NULL,
  PRIMARY KEY(memory_id, tag_id),
  FOREIGN KEY(memory_id) REFERENCES {{ .Namespace }}.memory(id) ON DELETE CASCADE,
  FOREIGN KEY(tag_id) REFERENCES {{ .Namespace }}.tag(id) ON DELETE CASCADE
);
