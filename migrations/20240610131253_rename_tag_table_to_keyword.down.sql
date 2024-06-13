ALTER TABLE memory_keyword RENAME CONSTRAINT memory_keyword_pkey TO memory_tag_pkey;
ALTER TABLE memory_keyword RENAME CONSTRAINT memory_keyword_keyword_id_fkey TO memory_tag_tag_id_fkey;
ALTER TABLE memory_keyword RENAME CONSTRAINT memory_keyword_memory_id_fkey TO memory_tag_memory_id_fkey;
ALTER TABLE memory_keyword RENAME COLUMN keyword_id TO tag_id;
ALTER TABLE memory_keyword RENAME TO memory_tag;

ALTER SEQUENCE keyword_id_seq RENAME TO tag_id_seq;
ALTER TABLE keyword RENAME CONSTRAINT keyword_pkey TO tag_pkey;
ALTER TABLE keyword RENAME CONSTRAINT keyword_name_key TO tag_name_key;
ALTER TABLE keyword RENAME TO tag;
