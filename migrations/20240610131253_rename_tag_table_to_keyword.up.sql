ALTER TABLE memory_tag RENAME COLUMN tag_id TO keyword_id;
ALTER TABLE memory_tag RENAME TO memory_keyword;
ALTER TABLE tag RENAME TO keyword;
