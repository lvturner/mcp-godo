# Quick Start: Category Feature

**Date**: 2025-10-25  
**Feature**: Category Feature for Todo Items  
**Status**: Draft

## Overview

The Category feature allows you to organize your todo items by grouping them into categories. This guide provides quick examples of how to use the category functionality through MCP tools.

## Basic Category Operations

### 1. Create a Category

Create a new category to organize your todos:

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

**Response**:
```
Category created: ID=1, Name=Work Tasks, Color=#3498db
```

### 2. List All Categories

View all available categories:

```json
{
  "tool": "get_all_categories",
  "arguments": {}
}
```

**Response**:
```
ID: 1, Name: Work Tasks, Description: Professional and work-related tasks, Color: #3498db, Created: 2025-10-25T08:30:00Z
ID: 2, Name: Personal, Description: Personal tasks and errands, Color: #e74c3c, Created: 2025-10-25T08:35:00Z
```

### 3. Get Category Details

Retrieve details for a specific category:

```json
{
  "tool": "get_category",
  "arguments": {
    "id": 1
  }
}
```

**Response**:
```
ID: 1, Name: Work Tasks, Description: Professional and work-related tasks, Color: #3498db, Created: 2025-10-25T08:30:00Z, Updated: 2025-10-25T08:30:00Z
```

## Working with Todos and Categories

### 4. Create a Todo in a Category

First create a todo, then assign it to a category:

```json
{
  "tool": "create_todo",
  "arguments": {
    "title": "Complete quarterly report",
    "description": "Prepare Q4 financial report for management"
  }
}
```

Then assign it to a category:

```json
{
  "tool": "assign_todo_to_category",
  "arguments": {
    "todo_id": 1,
    "category_id": 1
  }
}
```

**Response**:
```
Todo 'Complete quarterly report' assigned to category 'Work Tasks'
```

### 5. Get Todos in a Category

View all todos in a specific category:

```json
{
  "tool": "get_category_todos",
  "arguments": {
    "id": 1
  }
}
```

**Response**:
```
Todos in category 'Work Tasks':
- ID: 1, Title: Complete quarterly report, Status: Pending
- ID: 3, Title: Team meeting preparation, Status: Pending
```

### 6. Get Uncategorized Todos

View todos that haven't been assigned to any category:

```json
{
  "tool": "get_uncategorized_todos",
  "arguments": {}
}
```

**Response**:
```
Uncategorized todos:
- ID: 2, Title: Buy groceries, Status: Pending
- ID: 4, Title: Call dentist, Status: Pending
```

## Managing Categories

### 7. Update a Category

Modify category details:

```json
{
  "tool": "update_category",
  "arguments": {
    "id": 1,
    "name": "Work Projects",
    "color": "#2980b9"
  }
}
```

**Response**:
```
Category 1 updated successfully
```

### 8. Remove Todo from Category

Make a todo uncategorized:

```json
{
  "tool": "remove_todo_from_category",
  "arguments": {
    "todo_id": 1
  }
}
```

**Response**:
```
Todo removed from category successfully
```

### 9. Delete a Category

Remove a category (todos become uncategorized):

```json
{
  "tool": "delete_category",
  "arguments": {
    "id": 2
  }
}
```

**Response**:
```
Category 2 deleted successfully
```

## Common Workflows

### Workflow 1: Organize Existing Todos
1. Create categories for different areas of your life
2. Use `get_uncategorized_todos` to see unorganized items
3. Use `assign_todo_to_category` to organize them

### Workflow 2: Project-Based Organization
1. Create categories like "Work", "Personal", "Shopping"
2. Create todos and assign them to appropriate categories
3. Use `get_category_todos` to focus on one area at a time

### Workflow 3: Color-Coded Organization
1. Create categories with different colors
2. Use colors to quickly identify task types
3. Example: Red for urgent, blue for work, green for personal

## Best Practices

### Category Naming
- Use clear, descriptive names
- Keep names concise but meaningful
- Consider using consistent naming patterns

### Color Usage
- Use colors consistently across categories
- Choose colors that are visually distinct
- Consider color meanings (red=urgent, blue=calm, etc.)

### Organization Strategy
- Don't over-categorize - keep it simple
- Review uncategorized todos regularly
- Consider seasonal or project-based categories

## Error Handling

### Common Errors and Solutions

**Duplicate Category Name**:
```
Error: Category with name 'Work Tasks' already exists
```
**Solution**: Use a unique name or update the existing category

**Invalid Category ID**:
```
Error: Category not found
```
**Solution**: Check the category ID using `get_all_categories`

**Invalid Todo ID**:
```
Error: Todo not found
```
**Solution**: Verify the todo ID exists before assignment

## Tips and Tricks

1. **Start Simple**: Begin with 3-5 broad categories
2. **Regular Review**: Periodically check uncategorized todos
3. **Color Coding**: Use colors to create visual patterns
4. **Consistent Naming**: Develop a naming convention and stick to it
5. **Seasonal Categories**: Create temporary categories for specific projects or seasons

## Next Steps

- Explore advanced filtering options
- Consider integration with project feature (if available)
- Set up regular category maintenance routines
- Experiment with different organizational approaches