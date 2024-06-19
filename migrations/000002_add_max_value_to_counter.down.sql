-- Drop max_value column which is int default 0 from the counter table
ALTER TABLE counter DROP COLUMN IF EXISTS max_value;

