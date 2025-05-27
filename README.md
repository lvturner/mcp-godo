# Overview

A simple MCP server that exposes some basic functions for managing a todo list. Written in Go for my own learning purposes, and to fit my own workflows.
It will store your todos in a simple SQLLite database.

# Features

## 1. Add a Todo Item
**Tool:** `add_todo`  
**Parameters:**  
- `title` (required): The title of the todo item.  
- `due_date` (optional): Due date in ISO 8601 format (`2006-01-02T15:04:05Z`).  

## 2. Complete a Todo Item
**Tool:** `complete_todo`  
**Parameters:**  
- `id` (required): The ID of the todo item to mark as completed.

## 3. Uncomplete a Todo Item
**Tool:** `uncomplete_todo`  
**Parameters:**  
- `id` (required): The ID of the todo item to mark as uncompleted.

## 4. List All Todos
**Tool:** `list_todos`  
**Description:**  
Lists all todos with their IDs and completion status.  

## 5. Retrieve a Single Todo
**Tool:** `get_todo`  
**Parameters:**  
- `id` (required): The ID of the todo item to retrieve.

## 6. Delete a Todo Item
**Tool:** `delete_todo`  
**Parameters:**  
- `id` (required): The ID of the todo item to delete.

## 7. Get Active Todos
**Tool:** `get_active_todos`  
**Description:**  
Retrieves all todos that are not completed.

## 8. Get Completed Todos
**Tool:** `get_completed_todos`  
**Description:**  
Retrieves all todos that are completed.

## 9. Update a Todo's Due Date
**Tool:** `update_due_date`  
**Parameters:**  
- `id` (required): The ID of the todo item.  
- `due_date` (required): New due date in ISO 8601 format (`2006-01-02T15:04:05Z`).

## Example JSON configuration file
```json
{
  "name": "Todo list",
  "key": "todo-list",
  "description": "Basic todos",
  "command": "/path/to/compiled/binary",
  "env": {
    "STORAGE_TYPE": "sql",
    "DB_PATH": "/path/to/sqlite3/database.db"
  }
}
```
## Database schema
```sql
CREATE TABLE todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    completed INTEGER
, due_date DATETIME);
```

## Todos/Roadmap
- [ ] Improve DB Setup flow
- [ ] Tidy up main.go
- [ ] Consider adding more features like priority levels and tags
- [ ] Implement create_date field (and replace completed field with completion date) 
- [ ] Unit tests