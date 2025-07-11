-- migrations/0002_add_recurrence_support.sql
-- Adds support for recurring todos

BEGIN;

-- Add reference_id column to todos table
ALTER TABLE todos ADD COLUMN reference_id INT;

-- Create recurrence_patterns table
CREATE TABLE recurrence_patterns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    todo_id VARCHAR(255) NOT NULL,
    frequency VARCHAR(50) NOT NULL,
    `interval` INT NOT NULL,
    until DATETIME NULL,
    count INT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP()
) ENGINE=InnoDB;

-- Add index for better performance
CREATE INDEX idx_recurrence_patterns_todo_id ON recurrence_patterns(todo_id);

COMMIT;