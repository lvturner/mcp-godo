-- Drop the foreign key constraint first
ALTER TABLE todos DROP CONSTRAINT IF EXISTS fk_todos_category;

-- Drop the index
DROP INDEX IF EXISTS idx_todos_category_id;

-- Drop the category_id column
ALTER TABLE todos DROP COLUMN IF EXISTS category_id;