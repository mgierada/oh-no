-- Add max_value column which is int default 0 to counter table
ALTER TABLE ohno_counter ADD COLUMN max_value Int DEFAULT 0;
