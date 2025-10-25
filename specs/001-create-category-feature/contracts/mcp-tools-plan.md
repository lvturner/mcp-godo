# MCP Tools Plan: Category Feature

**Date**: 2025-10-25  
**Feature**: Category Feature for Todo Items  
**Status**: Draft

## Overview

This document outlines the MCP (Model Context Protocol) tools required for the Category feature. These tools will enable LLM agents to interact with the category management system through standardized MCP methods.

## MCP Tools

### 1. create_category
Creates a new category with the specified name, optional description, and optional color.

**Tool Name**: `create_category`  
**Description**: Create a new category for organizing todo items

**Parameters**:
- `name` (string, required): The name of the category
- `description` (string, optional): Description of the category
- `color` (string, optional): Hex color code for visual organization (e.g., "#FF5733")

**Returns**: 
- Success: Category details including ID, name, description, color, and creation timestamp
- Error: Descriptive error message if creation fails

**Example Usage**:
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

### 2. get_all_categories
Retrieves all categories in the system.

**Tool Name**: `get_all_categories`  
**Description**: Get a list of all available categories

**Parameters**: None

**Returns**: 
- Success: List of all categories with their details
- Empty list if no categories exist

**Example Usage**:
```json
{
  "tool": "get_all_categories",
  "arguments": {}
}
```

### 3. get_category
Retrieves a specific category by ID.

**Tool Name**: `get_category`  
**Description**: Get details of a specific category

**Parameters**:
- `id` (number, required): The ID of the category to retrieve

**Returns**: 
- Success: Category details
- Error: Not found message if category doesn't exist

**Example Usage**:
```json
{
  "tool": "get_category",
  "arguments": {
    "id": 1
  }
}
```

### 4. update_category
Updates an existing category's details.

**Tool Name**: `update_category`  
**Description**: Update category name, description, or color

**Parameters**:
- `id` (number, required): The ID of the category to update
- `name` (string, optional): New name for the category
- `description` (string, optional): New description for the category
- `color` (string, optional): New hex color code

**Returns**: 
- Success: Updated category details
- Error: Descriptive error message if update fails

**Example Usage**:
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

### 5. delete_category
Deletes a category from the system.

**Tool Name**: `delete_category`  
**Description**: Delete a category and make associated todos uncategorized

**Parameters**:
- `id` (number, required): The ID of the category to delete

**Returns**: 
- Success: Confirmation of deletion
- Error: Descriptive error message if deletion fails

**Example Usage**:
```json
{
  "tool": "delete_category",
  "arguments": {
    "id": 1
  }
}
```

### 6. get_category_todos
Retrieves all todo items associated with a specific category.

**Tool Name**: `get_category_todos`  
**Description**: Get all todos assigned to a specific category

**Parameters**:
- `id` (number, required): The ID of the category

**Returns**: 
- Success: List of todo items in the category
- Empty list if no todos are assigned to the category

**Example Usage**:
```json
{
  "tool": "get_category_todos",
  "arguments": {
    "id": 1
  }
}
```

### 7. assign_todo_to_category
Assigns an existing todo item to a category.

**Tool Name**: `assign_todo_to_category`  
**Description**: Assign a todo item to a specific category

**Parameters**:
- `todo_id` (number, required): The ID of the todo item to assign
- `category_id` (number, required): The ID of the category to assign the todo to

**Returns**: 
- Success: Confirmation of assignment
- Error: Descriptive error message if assignment fails

**Example Usage**:
```json
{
  "tool": "assign_todo_to_category",
  "arguments": {
    "todo_id": 1,
    "category_id": 1
  }
}
```

### 8. remove_todo_from_category
Removes a todo item from its category (makes it uncategorized).

**Tool Name**: `remove_todo_from_category`  
**Description**: Remove a todo item from its category

**Parameters**:
- `todo_id` (number, required): The ID of the todo item to remove from its category

**Returns**: 
- Success: Confirmation of removal
- Error: Descriptive error message if removal fails

**Example Usage**:
```json
{
  "tool": "remove_todo_from_category",
  "arguments": {
    "todo_id": 1
  }
}
```

### 9. get_uncategorized_todos
Retrieves all todo items that are not assigned to any category.

**Tool Name**: `get_uncategorized_todos`  
**Description**: Get all todo items without category assignment

**Parameters**: None

**Returns**: 
- Success: List of uncategorized todo items
- Empty list if all todos are categorized

**Example Usage**:
```json
{
  "tool": "get_uncategorized_todos",
  "arguments": {}
}
```

## Error Handling

All tools should implement consistent error handling:

1. **Validation Errors**: Return clear messages about missing or invalid parameters
2. **Not Found Errors**: Return appropriate messages when categories or todos don't exist
3. **Database Errors**: Return user-friendly messages while logging detailed errors
4. **Business Logic Errors**: Return clear explanations of business rule violations

## Response Format

### Success Response
```json
{
  "content": [
    {
      "type": "text",
      "text": "Category created: ID=1, Name=Work Tasks, Color=#3498db"
    }
  ]
}
```

### Error Response
```json
{
  "isError": true,
  "content": [
    {
      "type": "text",
      "text": "Error: Category name is required and must be a string"
    }
  ]
}
```

## Implementation Notes

1. **Consistency**: Follow the same patterns established in the project feature implementation
2. **Validation**: Implement proper input validation for all parameters
3. **Database Operations**: Use prepared statements and proper transaction handling
4. **Performance**: Consider adding pagination for list operations if categories/todos grow large
5. **Security**: Sanitize all inputs and implement proper access controls

## Testing Requirements

Each MCP tool should have corresponding tests:
- Unit tests for parameter validation
- Integration tests for database operations
- Contract tests for MCP protocol compliance
- Edge case testing for error conditions