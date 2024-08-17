CREATE TABLE IF NOT EXISTS {{ .Namespace }}.chunk(
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  memory_id uuid NOT NULL,
  sequence INTEGER NOT NULL,
  chunk TEXT NOT NULL,
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL,
  FOREIGN KEY(memory_id) REFERENCES {{ .Namespace }}.memory(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_chunk_memory_id ON {{ .Namespace }}.chunk(memory_id);

CREATE TABLE IF NOT EXISTS {{ .Namespace }}.chunk_keyword(
  keyword_id INTEGER NOT NULL,
  chunk_id uuid NOT NULL,
  PRIMARY KEY (keyword_id, chunk_id),
  FOREIGN KEY (chunk_id) REFERENCES {{ .Namespace }}.chunk(id) ON DELETE CASCADE,
  FOREIGN KEY (keyword_id) REFERENCES {{ .Namespace }}.keyword(id) ON DELETE CASCADE
);

comment on column {{ .Namespace }}.chunk_keyword.keyword_id is 'Chunk of text may not be associated with any keywords';
