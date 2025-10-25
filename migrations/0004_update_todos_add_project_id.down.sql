-- migrations/0004_update_todos_add_project_id.down.sql
-- Rolls back the project_id addition to todos table

BEGIN;

-- Drop the index first
DROP INDEX IF EXISTS idx_todos_project_id ON todos;

-- Remove the project_id column
ALTER TABLE todos DROP COLUMN IF EXISTS project_id;

COMMIT;