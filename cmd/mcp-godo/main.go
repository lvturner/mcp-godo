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
	handler := handler.NewHandler(todoService)

	// Add tool
	tool := mcp.NewTool("add_todo",
		mcp.WithDescription("Add a todo item to the list"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The title of the todo item"),
		),
		mcp.WithString("due_date",
			mcp.Description("The due date of the todo item in ISO 8601 format it should match the template '2006-01-02T15:04:05Z' if the user has not specified a time, assume midnight of the due date"),
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
}

