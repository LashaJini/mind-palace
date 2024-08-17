DROP TABLE IF EXISTS {{ .Namespace }}.chunk CASCADE;
DROP TABLE IF EXISTS {{ .Namespace }}.chunk_keyword;

DROP INDEX IF EXISTS {{ .Namespace }}.idx_chunk_memory_id;
