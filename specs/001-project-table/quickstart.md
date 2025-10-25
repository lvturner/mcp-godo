# Quick Start Guide: Project Table Feature

## Overview
This guide provides a quick overview of the Project Table feature implementation for the MCP Todo server. The feature adds project management capabilities to the existing todo system, allowing users to organize todos into projects.

## Key Components

### 1. Data Model
- **Projects**: Stores project information (id, name, description, timestamps)
- **Todos**: Enhanced with project_id field for project association
- **Relationships**: One-to-many relationship between projects and todos

### 2. MCP Tools
The following MCP tools will be available for project management:

#### Project Management Tools
- `create_project` - Create a new project
- `get_all_projects` - List all projects
- `get_project` - Get details of a specific project
- `update_project` - Update project information
- `delete_project` - Delete a project

#### Project-Todo Integration Tools
- `get_project_todos` - Get all todos for a specific project
- `add_todo_to_project` - Add a new todo to a project

### 3. Usage Examples

#### Creating a Project
```json
{
  "tool": "create_project",
  "arguments": {
    "name": "Website Redesign",
    "description": "Complete overhaul of company website"
  }
}
```

#### Adding a Todo to a Project
```json
{
  "tool": "add_todo_to_project",
  "arguments": {
    "title": "Design homepage mockup",
    "project_id": 1,
    "due_date": "2024-02-15T00:00:00Z"
  }
}
```

#### Getting Project Todos
```json
{
  "tool": "get_project_todos",
  "arguments": {
    "id": 1
  }
}
```

## Implementation Structure

### Files to be Created/Modified
1. **pkg/todo/project.go** - Project data structures and interface
2. **pkg/todo/project_mariadb.go** - Database implementation for projects
3. **pkg/handler/project_handler.go** - MCP tool handlers
4. **cmd/mcp-godo/main.go** - Add new MCP tools
5. **pkg/todo/todo.go** - Update TodoItem to include project_id
6. **pkg/todo/todo_mariadb.go** - Update to handle project_id

### Database Schema Changes
```sql
CREATE TABLE projects (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Add project_id to existing todos table
ALTER TABLE todos ADD COLUMN project_id INT DEFAULT NULL;
```

## Integration with Existing Features

### Todo Operations
- Existing todo operations remain unchanged
- New todos can be created with or without project association
- Project association is optional for backward compatibility

### Recurrence Patterns
- Recurrence patterns continue to work with project-associated todos
- Project association is preserved through recurrence cycles

## Error Handling

### Common Error Scenarios
- **Project not found**: Returns appropriate error when project ID doesn't exist
- **Duplicate project name**: Prevents creation of projects with duplicate names
- **Invalid project ID**: Validates project ID format and existence
- **Database errors**: Handles connection and query errors gracefully

## Testing Approach

### Unit Tests
- Test all new MCP tool handlers
- Test database operations for projects
- Test project-todo relationship operations
- Test error scenarios and edge cases

### Integration Tests
- Test MCP tool integration with existing todo system
- Test database schema changes
- Test project-todo associations

## Migration Path

### For Existing Installations
1. Run database migration to add projects table
2. Run database migration to add project_id to todos table
3. Deploy new MCP server with project tools
4. Existing todo operations remain unchanged

### For New Installations
1. Fresh database schema includes projects table and project_id field
2. All project and todo tools available immediately

## Performance Considerations

### Database Optimization
- Index on projects.name for unique constraint
- Index on todos.project_id for project-todo queries
- Proper foreign key relationships for data integrity

### MCP Tool Performance
- Efficient database queries for project operations
- Minimal overhead for existing todo operations
- Optimized project-todo relationship queries

## Security Considerations

### Data Validation
- Input validation for all project fields
- SQL injection prevention through prepared statements
- Proper error handling to prevent information leakage

### Access Control
- Consider implementing project-level permissions if needed
- Ensure proper authorization for project operations
- Validate user access to project data

This quick start guide provides the foundation for implementing the Project Table feature while maintaining compatibility with the existing todo system.