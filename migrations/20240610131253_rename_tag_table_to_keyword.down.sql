ALTER TABLE memory_keyword RENAME COLUMN keyword_id TO tag_id;
ALTER TABLE memory_keyword RENAME TO memory_tag;
ALTER TABLE keyword RENAME TO tag;
