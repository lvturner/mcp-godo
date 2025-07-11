-- migrations/0002_add_recurrence_support.down.sql
-- Rolls back the recurrence pattern support

BEGIN;

-- Drop the index first
DROP INDEX IF EXISTS idx_recurrence_patterns_todo_id ON recurrence_patterns;

-- Drop the recurrence_patterns table
DROP TABLE IF EXISTS recurrence_patterns;

-- Remove the reference_id column
ALTER TABLE todos DROP COLUMN IF EXISTS reference_id;

COMMIT;