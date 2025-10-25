-- Drop the update trigger
DROP TRIGGER IF EXISTS update_categories_timestamp;

-- Drop the index
DROP INDEX IF EXISTS idx_categories_name;

-- Drop the categories table
DROP TABLE IF EXISTS categories;