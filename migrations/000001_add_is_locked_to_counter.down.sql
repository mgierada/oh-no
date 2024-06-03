-- Drop is_locked column which is bool default false to counter table
ALTER TABLE counter DROP COLUMN IF EXISTS is_locked;

