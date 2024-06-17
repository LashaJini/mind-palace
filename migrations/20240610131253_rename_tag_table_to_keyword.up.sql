ALTER TABLE {{ .Namespace }}.memory_tag RENAME CONSTRAINT memory_tag_pkey TO memory_keyword_pkey;
ALTER TABLE {{ .Namespace }}.memory_tag RENAME CONSTRAINT memory_tag_tag_id_fkey TO memory_keyword_keyword_id_fkey;
ALTER TABLE {{ .Namespace }}.memory_tag RENAME CONSTRAINT memory_tag_memory_id_fkey TO memory_keyword_memory_id_fkey;
ALTER TABLE {{ .Namespace }}.memory_tag RENAME COLUMN tag_id TO keyword_id;
ALTER TABLE {{ .Namespace }}.memory_tag RENAME TO memory_keyword;

ALTER SEQUENCE {{ .Namespace }}.tag_id_seq RENAME TO keyword_id_seq;
ALTER TABLE {{ .Namespace }}.tag RENAME CONSTRAINT tag_pkey TO keyword_pkey;
ALTER TABLE {{ .Namespace }}.tag RENAME CONSTRAINT tag_name_key TO keyword_name_key;
ALTER TABLE {{ .Namespace }}.tag RENAME TO keyword;
