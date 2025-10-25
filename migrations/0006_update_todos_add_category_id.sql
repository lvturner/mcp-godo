-- Add category_id column to todos table for category relationship
ALTER TABLE todos ADD COLUMN category_id INTEGER DEFAULT NULL;

-- Add foreign key constraint with ON DELETE SET NULL
-- Note: SQLite supports foreign keys but they must be enabled per connection
ALTER TABLE todos 
ADD CONSTRAINT fk_todos_category 
FOREIGN KEY (category_id) REFERENCES categories(id) 
ON DELETE SET NULL;

-- Create index for performance on category_id field
CREATE INDEX idx_todos_category_id ON todos(category_id);