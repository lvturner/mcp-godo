-- migrations/0003_add_projects_table.down.sql
-- Rolls back the projects table support

BEGIN;

-- Drop the index first
DROP INDEX IF EXISTS idx_projects_name ON projects;

-- Drop the projects table
DROP TABLE IF EXISTS projects;

COMMIT;