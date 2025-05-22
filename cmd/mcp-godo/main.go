package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var todoService todo.TodoService
var config todo.Config

func loadConfig(){
	config.StorageType = os.Getenv("STORAGE_TYPE")
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
		mcp.WithDescription("Add a new todo item, always call list_todos before calling this to help avoid adding duplicate items, always remember to actually call this function to add the todo item correctly to the list"),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The title of the todo item"),
		),
	)

	// Add tool handler
	s.AddTool(tool, addTodoHandler)
	

	completeTodoTool := mcp.NewTool("complete_todo",
		mcp.WithDescription("Complete a single todo item by ID, remember you can list todos in order to check for an appropriate id"),
		mcp.WithString("id",
			mcp.Required(),
			mcp.Description("The ID of the todo item"),
		),
	)
	s.AddTool(completeTodoTool, completeTodoHandler)
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
	
	todoService.AddTodo(title)
	
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