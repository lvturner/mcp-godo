# Data Model: Category Feature

**Date**: 2025-10-25  
**Feature**: Category Feature for Todo Items  
**Status**: Draft

## Entities

### Category
Represents a grouping of todo items by type or theme.

**Attributes**:
- `id` (int64): Primary key, auto-incrementing unique identifier
- `name` (string): Unique category name, required field
- `description` (*string): Optional category description, can be NULL
- `color` (*string): Optional hex color code for visual organization (e.g., "#FF5733")
- `created_at` (time.Time): Timestamp when category was created
- `updated_at` (time.Time): Timestamp when category was last updated

**Validation Rules**:
- Name must be unique across all categories
- Name cannot be empty
- Name should have reasonable length limit (255 characters)
- Color must be valid hex format if provided (optional validation)
- Description can be empty or NULL

**State Transitions**:
- **Create**: New category created with generated ID and timestamps
- **Update**: Name, description, and/or color can be modified, updated_at timestamp changes
- **Delete**: Category removed from system, associated todos become uncategorized

### TodoItem (Updated)
Existing todo entity with added category relationship.

**Additional Attributes**:
- `category_id` (*int64): Foreign key to categories table, optional (can be NULL)

**Validation Rules**:
- category_id must reference an existing category if provided
- A todo can have either a project_id OR a category_id, or neither, but not both (business rule to be determined)

**Relationships**:
- Many-to-one with Category (many todos can belong to one category)
- Maintains existing relationships with Project (if implemented)

## Database Schema

### Categories Table
```sql
CREATE TABLE categories (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    color VARCHAR(7), -- Hex color code format: #RRGGBB
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Indexes for performance
CREATE INDEX idx_categories_name ON categories(name);
```

### Todos Table Update
```sql
-- Add category_id to existing todos table
ALTER TABLE todos ADD COLUMN category_id BIGINT DEFAULT NULL;

-- Add foreign key constraint
ALTER TABLE todos 
ADD CONSTRAINT fk_todos_category 
FOREIGN KEY (category_id) REFERENCES categories(id) 
ON DELETE SET NULL;

-- Index for performance
CREATE INDEX idx_todos_category_id ON todos(category_id);
```

## Entity Relationships

### Category → TodoItem
- **Type**: One-to-Many
- **Cardinality**: One category can have many todo items
- **Optional**: Yes (todos can exist without categories)
- **Delete Rule**: SET NULL (when category is deleted, todos become uncategorized)

### TodoItem → Category
- **Type**: Many-to-One
- **Cardinality**: Many todo items can belong to one category
- **Optional**: Yes (todo items can exist without categories)
- **Constraint**: Foreign key constraint ensures referential integrity

## Business Rules

1. **Unique Category Names**: Category names must be unique within the system
2. **Optional Categorization**: Todo items can be created without specifying a category
3. **Single Categorization**: A todo item can belong to at most one category (no many-to-many relationship)
4. **Category Deletion**: When a category is deleted, associated todos become uncategorized (category_id set to NULL)
5. **Color Validation**: If provided, color must be a valid hex color code format
6. **Name Length**: Category names are limited to 255 characters

## Data Validation

### Category Creation/Update
- Name: Required, unique, max 255 characters
- Description: Optional, no length limit
- Color: Optional, must be valid hex format (#RRGGBB)

### Todo Creation/Update with Category
- category_id: Optional, must reference existing category if provided
- Cannot assign to both project and category simultaneously (business rule)

## Migration Strategy

### Phase 1: Create Categories Table
- Create categories table with proper schema
- Add indexes for performance
- No impact on existing data

### Phase 2: Update Todos Table
- Add category_id column to todos table
- Add foreign key constraint
- Add index for performance
- Existing todos will have category_id = NULL (uncategorized)

## Query Patterns

### Common Queries
1. **Get all categories**: `SELECT * FROM categories ORDER BY created_at DESC`
2. **Get todos by category**: `SELECT * FROM todos WHERE category_id = ?`
3. **Get uncategorized todos**: `SELECT * FROM todos WHERE category_id IS NULL`
4. **Get category with todo count**: `SELECT c.*, COUNT(t.id) as todo_count FROM categories c LEFT JOIN todos t ON c.id = t.category_id GROUP BY c.id`

### Performance Considerations
- Indexes on category_id and name fields optimize common queries
- Foreign key constraints ensure data integrity
- LEFT JOIN for counting todos prevents missing empty categories

## Future Considerations

1. **Multi-user Support**: Consider adding user_id to categories if multi-user functionality is added
2. **Category Hierarchy**: Could support nested categories with parent_id field
3. **Category Templates**: Could support predefined category sets
4. **Bulk Operations**: Consider bulk category assignment for multiple todos