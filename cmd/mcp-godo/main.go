package main

import (
	"fmt"
	"os"

	"mcp-godo/pkg/handler"
	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var todoService todo.TodoService
var projectService todo.ProjectService
var config todo.Config

func loadConfig(){
	config.StorageType = os.Getenv("STORAGE_TYPE")
	config.SQLDBPath = os.Getenv("DB_PATH")
	config.HTTPPort = os.Getenv("HTTP_PORT")
}

func main() {
	var err error
	loadConfig()
	
	todoService, err = todo.NewTodoServiceFromConfig(config)	
	if err != nil {
		fmt.Println("Error creating todo service:", err)
		return
	}

	// Initialize project service using the same database connection
	projectService, err = todo.NewProjectServiceFromConfig(config)
	if err != nil {
		fmt.Println("Error creating project service:", err)
		return
	}

	// Create a new MCP server
	s := server.NewMCPServer(
		"Todo MCP",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	addTools(s)

	httpServer := server.NewStreamableHTTPServer(s)
	if err := httpServer.Start(fmt.Sprintf(":%s", config.HTTPPort)); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func addTools(s *server.MCPServer) {
	handler := handler.NewHandlerWithProject(todoService, projectService)

	// Add tool with project_id support
	tool := mcp.NewTool("add_todo",
		mcp.WithDescription("Add a todo item to the list"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The title of the todo item"),
		),
		mcp.WithString("due_date",
			mcp.Description("The due date of the todo item in ISO 8601 format it should match the template '2006-01-02T15:04:05Z' if the user has not specified a time, assume midnight of the due date"),
		),
		mcp.WithNumber("project_id",
			mcp.Description("The ID of the project to assign this todo to (optional)"),
		),
	)
	s.AddTool(tool, handler.AddTodoHandler)
	
	completeTodoTool := mcp.NewTool("complete_todo",
		mcp.WithDescription("Complete a single todo item by ID - you may need to call get_active_todos or list_todos in order to get the correct ID - lookup by title or other attributes won't work with this call"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(completeTodoTool, handler.CompleteTodoHandler)
	
	unCompleteTodoTool := mcp.NewTool("uncomplete_todo",
		mcp.WithDescription("Uncomplete a single todo item by ID (mark it as undone) you may need to call get_active_todos or list_todos in order to get the correct ID - lookup by title or other attributes won't work with this call"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(unCompleteTodoTool, handler.UnCompleteTodoHandler)
	
	listTodosTool := mcp.NewTool("list_todos",
		mcp.WithDescription("Lists all todo items with their IDs and completion status."),
	)
	s.AddTool(listTodosTool, handler.ListTodosHandler)
	
	getTodoTool := mcp.NewTool("get_todo",
		mcp.WithDescription("Retrieve details of a single todo item by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(getTodoTool, handler.GetTodoHandler)
	
	deleteTodoTool := mcp.NewTool("delete_todo",
		mcp.WithDescription("Delete a single todo item by ID"),		
		mcp.WithString("id",
			mcp.Required(),			
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(deleteTodoTool, handler.DeleteTodoHandler)
	
	getActiveTodosTool := mcp.NewTool("get_active_todos",
		mcp.WithDescription("Retrieve all active (not completed) todos (generally prefer this over list_todos)"),
	)
	s.AddTool(getActiveTodosTool, handler.GetActiveTodosHandler)
	
	getCompletedTodosTool := mcp.NewTool("get_completed_todos",
		mcp.WithDescription("Retrieve all completed todos"),
	)
	s.AddTool(getCompletedTodosTool, handler.GetCompletedTodosHandler)
	
	updateDueDateTool := mcp.NewTool("update_due_date",
		mcp.WithDescription("Update the due date of a single todo item by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
		mcp.WithString("due_date",
			mcp.Required(),
			mcp.Description("The new due date for the todo item"),
		),
	)
	s.AddTool(updateDueDateTool, handler.UpdateDueDateHandler)

	titleSearchTool := mcp.NewTool("title_search",
		mcp.WithDescription("Search todos by title, if this returns nothing or an error, call get_active_todos to find the todo "),
		mcp.WithString("query",
			mcp.Required(),
			mcp.Description("The search query to use"),
		),
		mcp.WithString("active_only",
			mcp.DefaultBool(true),
			mcp.Description("Whether to only search active todos or not"),
		),
	)
	s.AddTool(titleSearchTool, handler.TitleSearchHandler)

	// Add recurrence pattern tool
	addRecurrencePatternTool := mcp.NewTool("add_recurrence_pattern",
		mcp.WithDescription("Add a recurrence pattern to a todo item"),
		mcp.WithString("todo_id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
		mcp.WithString("frequency",
			mcp.Required(),
			mcp.Description("The frequency of the recurrence (e.g., 'daily', 'weekly', 'monthly', 'yearly')"),
		),
		mcp.WithNumber("interval",
			mcp.Required(),
			mcp.Description("The interval between recurrences (e.g., 1 for every day/week/month/year)"),
		),
		mcp.WithString("until",
			mcp.Description("The end date for the recurrence pattern in ISO 8601 format (optional)"),
		),
		mcp.WithNumber("count",
			mcp.Description("The number of times the recurrence should occur (optional)"),
		),
	)
	s.AddTool(addRecurrencePatternTool, handler.AddRecurrencePatternHandler)

	// Get recurrence pattern tool
	getRecurrencePatternTool := mcp.NewTool("get_recurrence_pattern",
		mcp.WithDescription("Retrieve a recurrence pattern by its ID"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("The ID of the recurrence pattern"),
		),
	)
	s.AddTool(getRecurrencePatternTool, handler.GetRecurrencePatternHandler)

	// Add project management tools
	addProjectTools(s, handler)
}

func addProjectTools(s *server.MCPServer, handler *handler.Handler) {
	// Create project tool
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

	// Get all projects tool
	getAllProjectsTool := mcp.NewTool("get_all_projects",
		mcp.WithDescription("Retrieve all projects"),
	)
	s.AddTool(getAllProjectsTool, handler.GetAllProjectsHandler)

	// Get project tool
	getProjectTool := mcp.NewTool("get_project",
		mcp.WithDescription("Retrieve details of a single project by ID"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("The ID of the project"),
		),
	)
	s.AddTool(getProjectTool, handler.GetProjectHandler)

	// Update project tool
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

	// Delete project tool
	deleteProjectTool := mcp.NewTool("delete_project",
		mcp.WithDescription("Delete a project by ID"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("The ID of the project"),
		),
	)
	s.AddTool(deleteProjectTool, handler.DeleteProjectHandler)

	// Get project todos tool
	getProjectTodosTool := mcp.NewTool("get_project_todos",
		mcp.WithDescription("Retrieve all todos for a specific project"),
		mcp.WithNumber("id",
			mcp.Required(),
			mcp.Description("The ID of the project"),
		),
	)
	s.AddTool(getProjectTodosTool, handler.GetProjectTodosHandler)

	// Add todo to project tool
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
}
