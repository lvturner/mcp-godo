package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var todoService todo.TodoService
var config todo.Config

func loadConfig(){
	config.StorageType = os.Getenv("STORAGE_TYPE")
	config.SQLDBPath = os.Getenv("DB_PATH")
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
	addResources(s)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Not really in use as my client doesn't yet support resources properly, here for future expansion. *untested*
func addResources(s *server.MCPServer) {
	listTodosTemplate := mcp.NewResource(
		"todos",
		"List of todos",
		mcp.WithResourceDescription("Full list of todos"),
		mcp.WithMIMEType("application/json"),
	)
	
	s.AddResource(listTodosTemplate, listTodosResourceHandler)

	
	singleTodoTemplate := mcp.NewResource(
		"todos://{id}",
		"Single todo list item",
		mcp.WithResourceDescription("Single todo list item"),
		mcp.WithMIMEType("application/json"),
	)	
	
	s.AddResource(singleTodoTemplate, getSingleTodoResourceHandler)
}

func addTools(s *server.MCPServer) {
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

	// Add tool handler
	s.AddTool(tool, addTodoHandler)
	
	completeTodoTool := mcp.NewTool("complete_todo",
		mcp.WithDescription("Complete a single todo item by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(completeTodoTool, completeTodoHandler)
	
	unCompleteTodoTool := mcp.NewTool("uncomplete_todo",
		mcp.WithDescription("Uncomplete a single todo item by ID (mark it as undone)"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(unCompleteTodoTool, unCompleteTodoHandler)
	
	listTodosTool := mcp.NewTool("list_todos",
		mcp.WithDescription("Lists all todo items with their IDs and completion status."),
	)
	s.AddTool(listTodosTool, listTodosHandler)
	
	getTodoTool := mcp.NewTool("get_todo",
		mcp.WithDescription("Retrieve details of a single todo item by ID"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(getTodoTool, getTodoHandler)
	
	deleteTodoTool := mcp.NewTool("delete_todo",
		mcp.WithDescription("Delete a single todo item by ID"),		
		mcp.WithString("id",
			mcp.Required(),			
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(deleteTodoTool, deleteTodoHandler)
	
	getActiveTodosTool := mcp.NewTool("get_active_todos",
		mcp.WithDescription("Retrieve all active (not completed) todos"),
	)
	s.AddTool(getActiveTodosTool, getActiveTodosHandler)
	
	getCompletedTodosTool := mcp.NewTool("get_completed_todos",
		mcp.WithDescription("Retrieve all completed todos"),
	)
	s.AddTool(getCompletedTodosTool, getCompletedTodosHandler)
	
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
	s.AddTool(updateDueDateTool, updateDueDateHandler)
}

func updateDueDateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok{
		return nil, fmt.Errorf("invalid id")
	}
	dueDateStr, ok := request.GetArguments()["due_date"].(string)
	if !ok{
		return nil, fmt.Errorf("invalid due date")
	}
	dueDate, err := time.Parse(time.RFC3339, dueDateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse due date: %w", err)
	}
	todo, err := todoService.SetDueDate(id, dueDate)
	if err != nil{
		return nil, fmt.Errorf("failed to update due date: %w", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Todo updated: ID=%s, Title=%s, Due Date=%s", todo.ID, todo.Title, todo.DueDate.Format(time.RFC3339))), nil
}

func unCompleteTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok{
		return nil, fmt.Errorf("invalid id")
	}
	todo, err := todoService.UnCompleteTodo(id)
	if err != nil{
		return nil, fmt.Errorf("failed to uncomplete todo: %w", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Todo uncompleted: ID=%s, Title=%s", todo.ID, todo.Title)), nil
}

func getCompletedTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todos := todoService.GetCompletedTodos()
	if len(todos) == 0 {
		return mcp.NewToolResultText("No completed todos found"), nil
	}
	var resultText string
	for _, todo := range todos {
		status := "Completed"
		resultText += fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s\n", todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate)
	}
	return mcp.NewToolResultText(resultText), nil
}

func getActiveTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todos := todoService.GetActiveTodos()
	if len(todos) == 0 {
		return mcp.NewToolResultText("No active todos found"), nil
	}
	var resultText string
	for _, todo := range todos {
		status := "Incomplete"
		resultText += fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s\n", todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate)
	}
	return mcp.NewToolResultText(resultText), nil
}

func deleteTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok {
		return nil, errors.New("id must be a string")
	}
	todo, err := todoService.DeleteTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete todo: %w", err)	
	}

	resultText := fmt.Sprintf("Deleted Todo: ID=%s, Title=%s", todo.ID, todo.Title)
	return mcp.NewToolResultText(resultText), nil
}

func getTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok {
		return nil, errors.New("id must be a string")
	}
	todo, err := todoService.GetTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	status := "Incomplete"
	if todo.Completed {
		status = "Complete"
	}

	resultText := fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s\n", todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate)
	
	return mcp.NewToolResultText(resultText), nil
}


func listTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todos := todoService.GetAllTodos()
	var todosText []string
	for _, todo := range todos {
		status := "Incomplete"
		if todo.Completed {
			status = "Complete"
		}
		todosText = append(todosText, fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s\n", todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate))
	}
	return mcp.NewToolResultText(strings.Join(todosText, "\n")), nil
}

func completeTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok {
		return nil, errors.New("id must be a string")
	}

	completedTodo, err := todoService.CompleteTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to complete todo: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Todo %s completed", completedTodo.Title)), nil
}

func addTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := request.GetArguments()["title"].(string)
	if !ok {
		return nil, errors.New("title must be a string")
	}
	dueDateRaw, ok := request.GetArguments()["due_date"]
	if ok {
		dueDateStr, ok := dueDateRaw.(string)
		if !ok {
			return nil, errors.New("due_date must be a string")
		}
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, err
		}

		todoService.AddTodo(title, &dueDate)
	} else {
		todoService.AddTodo(title, nil)
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("%s added to todo list", title)), nil
}

func listTodosResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	todos := todoService.GetAllTodos()
	
	jsonData, err := json.Marshal(todos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal todos to JSON: %w", err)
	}
	
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI: request.Params.URI,
			MIMEType: "application/json",
			Text: string(jsonData),
		},
	}, nil
}

func getSingleTodoResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	id := extractIDFromURI(request.Params.URI)
	todo, err := todoService.GetTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo item: %w", err)
	}	
	
	jsonData, err := json.Marshal(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal todo item: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:	request.Params.URI,
			MIMEType: "application/json",
			Text: string(jsonData),
		},
	}, nil
}


func extractIDFromURI(uri string) string {
    parsed, err := url.Parse(uri)
    if err != nil || parsed.Scheme != "todos" || parsed.Opaque == "" {
        return ""
    }

    return strings.TrimPrefix(parsed.Path, "/") 
}