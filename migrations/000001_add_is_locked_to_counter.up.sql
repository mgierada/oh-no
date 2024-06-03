-- Add is_locked column which is bool default false to counter table
ALTER TABLE counter ADD COLUMN is_locked BOOLEAN DEFAULT FALSE;
