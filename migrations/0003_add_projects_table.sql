-- migrations/0003_add_projects_table.sql
-- Adds support for project management

BEGIN;

-- Create projects table
CREATE TABLE projects (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
) ENGINE=InnoDB;

-- Add index for better performance on name lookups
CREATE INDEX idx_projects_name ON projects(name);

COMMIT;