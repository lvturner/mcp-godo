-- migrations/0004_update_todos_add_project_id.sql
-- Adds project_id foreign key to todos table

BEGIN;

-- Add project_id column to todos table
ALTER TABLE todos ADD COLUMN project_id INT DEFAULT NULL;

-- Add index for better performance on project_id lookups
CREATE INDEX idx_todos_project_id ON todos(project_id);

COMMIT;