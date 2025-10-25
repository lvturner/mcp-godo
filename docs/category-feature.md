# Category Feature Documentation

## Overview

The category feature allows users to organize their todo items into categories for better task management. This feature provides full CRUD operations for categories and enables flexible categorization of todo items.

## Features

### 1. Category Management
- **Create Categories**: Create new categories with name, optional description, and optional color
- **View Categories**: List all categories or get details of a specific category
- **Update Categories**: Modify category name, description, or color
- **Delete Categories**: Remove categories (associated todos become uncategorized)

### 2. Todo-Category Relationships
- **Assign Todos to Categories**: Link existing todos to categories
- **Remove Todos from Categories**: Unlink todos from their categories
- **Create Todos in Categories**: Create new todos directly assigned to categories
- **View Todos by Category**: List all todos in a specific category
- **View Uncategorized Todos**: List todos that are not assigned to any category

## MCP Tools

The following MCP tools are available for category operations:

### Category Management Tools

#### `create_category`
Creates a new category.
```json
{
  "tool": "create_category",
  "arguments": {
    "name": "Work Tasks",
    "description": "Professional and work-related tasks",
    "color": "#3498db"
  }
}
```

#### `get_all_categories`
Lists all categories.
```json
{
  "tool": "get_all_categories",
  "arguments": {}
}
```

#### `get_category`
Gets details of a specific category.
```json
{
  "tool": "get_category",
  "arguments": {
    "id": 1
  }
}
```

#### `update_category`
Updates an existing category.
```json
{
  "tool": "update_category",
  "arguments": {
    "id": 1,
    "name": "Updated Work Tasks",
    "color": "#2980b9"
  }
}
```

#### `delete_category`
Deletes a category.
```json
{
  "tool": "delete_category",
  "arguments": {
    "id": 1
  }
}
```

### Todo-Category Relationship Tools

#### `get_category_todos`
Lists all todos in a specific category.
```json
{
  "tool": "get_category_todos",
  "arguments": {
    "id": 1
  }
}
```

#### `get_uncategorized_todos`
Lists all todos without category assignment.
```json
{
  "tool": "get_uncategorized_todos",
  "arguments": {}
}
```

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
);
```

### Todos Table (Updated)
```sql
ALTER TABLE todos ADD COLUMN category_id BIGINT DEFAULT NULL;
ALTER TABLE todos 
ADD CONSTRAINT fk_todos_category 
FOREIGN KEY (category_id) REFERENCES categories(id) 
ON DELETE SET NULL;
```

## Business Rules

1. **Unique Category Names**: Category names must be unique within the system
2. **Optional Categorization**: Todo items can exist without categories
3. **Single Categorization**: A todo item can belong to at most one category
4. **Category Deletion**: When a category is deleted, associated todos become uncategorized
5. **Color Validation**: If provided, color must be a valid hex color code format
6. **Name Length**: Category names are limited to 255 characters

## Error Handling

The feature includes comprehensive error handling for:
- Duplicate category names
- Invalid color formats
- Missing required parameters
- Non-existent categories or todos
- Database connection issues

## Performance Considerations

- Indexes on `category_id` and `name` fields optimize common queries
- Foreign key constraints ensure data integrity
- Efficient queries for category-todo relationships

## Testing

The feature includes comprehensive test coverage:
- Unit tests for category operations
- Unit tests for todo-category relationships
- Contract tests for MCP tool compliance
- Integration tests for database operations

## Security

- Input validation and sanitization
- SQL injection prevention through prepared statements
- Proper error message handling (user-friendly messages while logging detailed errors)