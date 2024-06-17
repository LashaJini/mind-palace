CREATE TABLE IF NOT EXISTS {{ .Namespace }}.memory(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS {{ .Namespace }}.resource(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  memory_id uuid NOT NULL,
  file_path TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  FOREIGN KEY(memory_id) REFERENCES {{ .Namespace }}.memory(id) ON DELETE CASCADE
);
