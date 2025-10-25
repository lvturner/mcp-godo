# MCP Tools Plan for Project Table Feature

Based on the existing MCP server implementation in `cmd/mcp-godo/main.go`, this document outlines the plan for adding project management tools to the MCP server.

## Project Data Structures

Following the pattern of the existing `TodoItem` and `TodoService` in `pkg/todo/todo.go`, we need to define:

### Project Struct
```go
type Project struct {
    ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description *string   `json:"description"` // pointer to handle NULL in database
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### ProjectService Interface
```go
type ProjectService interface {
    CreateProject(name string, description *string) (Project, error)
    GetAllProjects() []Project
    GetProject(id int64) (Project, error)
    UpdateProject(id int64, name string, description *string) (Project, error)
    DeleteProject(id int64) (Project, error)
    GetProjectTodos(id int64) []TodoItem
}
```

## MCP Tools to Add

Following the pattern of existing todo tools in `cmd/mcp-godo/main.go`, add these tools to the `addTools` function:

### 1. create_project
```go
createProjectTool := mcp.NewTool("create_project",
    mcp.WithDescription("Create a new project"),
    mcp.WithString("name",
        mcp.Required(),
        mcp.Description("The name of the project"),
    ),
    mcp.WithString("description",
        mcp.Description("The description of the project (optional)"),
    ),
)
s.AddTool(createProjectTool, handler.CreateProjectHandler)
```

### 2. get_all_projects
```go
getAllProjectsTool := mcp.NewTool("get_all_projects",
    mcp.WithDescription("Retrieve all projects"),
)
s.AddTool(getAllProjectsTool, handler.GetAllProjectsHandler)
```

### 3. get_project
```go
getProjectTool := mcp.NewTool("get_project",
    mcp.WithDescription("Retrieve details of a single project by ID"),
    mcp.WithNumber("id",
        mcp.Required(),
        mcp.Description("The ID of the project"),
    ),
)
s.AddTool(getProjectTool, handler.GetProjectHandler)
```

### 4. update_project
```go
updateProjectTool := mcp.NewTool("update_project",
    mcp.WithDescription("Update a project by ID"),
    mcp.WithNumber("id",
        mcp.Required(),
        mcp.Description("The ID of the project"),
    ),
    mcp.WithString("name",
        mcp.Description("The new name of the project"),
    ),
    mcp.WithString("description",
        mcp.Description("The new description of the project"),
    ),
)
s.AddTool(updateProjectTool, handler.UpdateProjectHandler)
```

### 5. delete_project
```go
deleteProjectTool := mcp.NewTool("delete_project",
    mcp.WithDescription("Delete a project by ID"),
    mcp.WithNumber("id",
        mcp.Required(),
        mcp.Description("The ID of the project"),
    ),
)
s.AddTool(deleteProjectTool, handler.DeleteProjectHandler)
```

### 6. get_project_todos
```go
getProjectTodosTool := mcp.NewTool("get_project_todos",
    mcp.WithDescription("Retrieve all todos for a specific project"),
    mcp.WithNumber("id",
        mcp.Required(),
        mcp.Description("The ID of the project"),
    ),
)
s.AddTool(getProjectTodosTool, handler.GetProjectTodosHandler)
```

### 7. add_todo_to_project
```go
addTodoToProjectTool := mcp.NewTool("add_todo_to_project",
    mcp.WithDescription("Add a todo item to a specific project"),
    mcp.WithString("title",
        mcp.Required(),
        mcp.Description("The title of the todo item"),
    ),
    mcp.WithNumber("project_id",
        mcp.Required(),
        mcp.Description("The ID of the project to add the todo to"),
    ),
    mcp.WithString("due_date",
        mcp.Description("The due date of the todo item in ISO 8601 format"),
    ),
)
s.AddTool(addTodoToProjectTool, handler.AddTodoToProjectHandler)
```

## Handler Implementation Plan

Create a new handler file `pkg/handler/project_handler.go` with methods:

1. `CreateProjectHandler(ctx context.Context, params map[string]interface{}) (*mcp.CallToolResult, error)`
2. `GetAllProjectsHandler(ctx context.Context, params map[string]interface{}) (*mcp.CallToolResult, error)`
3. `GetProjectHandler(ctx context.Context, params map[string]interface{}) (*mcp.CallToolResult, error)`
4. `UpdateProjectHandler(ctx context.Context, params map[string]interface{}) (*mcp.CallToolResult, error)`
5. `DeleteProjectHandler(ctx context.Context, params map[string]interface{}) (*mcp.CallToolResult, error)`
6. `GetProjectTodosHandler(ctx context.Context, params map[string]interface{}) (*mcp.CallToolResult, error)`
7. `AddTodoToProjectHandler(ctx context.Context, params map[string]interface{}) (*mcp.CallToolResult, error)`

## Database Implementation Plan

Create `pkg/todo/project_mariadb.go` following the pattern of `todo_mariadb.go` with methods:

1. `CreateProject(name string, description *string) (Project, error)`
2. `GetAllProjects() []Project`
3. `GetProject(id int64) (Project, error)`
4. `UpdateProject(id int64, name string, description *string) (Project, error)`
5. `DeleteProject(id int64) (Project, error)`
6. `GetProjectTodos(id int64) []TodoItem`

## Integration Points

1. Update `pkg/todo/todo.go` to include project_id in TodoItem struct
2. Update `pkg/todo/todo_mariadb.go` methods to handle project_id field
3. Update existing todo creation methods to support optional project_id parameter
4. Ensure foreign key constraints are properly handled in database operations

## Error Handling

Follow the same error handling patterns as existing todo tools:
- Return descriptive error messages
- Handle database connection errors
- Validate input parameters
- Handle not-found scenarios appropriately

## Testing Considerations

Plan for unit tests covering:
- All new MCP tool handlers
- Database operations for projects
- Project-todo relationship operations
- Error scenarios and edge cases