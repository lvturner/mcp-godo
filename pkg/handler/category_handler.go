package handler

import (
	"context"
	"fmt"

	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
)

// CategoryHandler handles MCP tools for category operations
type CategoryHandler struct {
	categoryService todo.CategoryService
	todoService     todo.TodoService
}

// NewCategoryHandler creates a new category handler
func NewCategoryHandler(categoryService todo.CategoryService, todoService todo.TodoService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		todoService:     todoService,
	}
}

// CreateCategoryHandler handles the create_category MCP tool
func (h *CategoryHandler) CreateCategoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	name, ok := request.GetArguments()["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required and must be a string")
	}
	
	// Optional parameters
	var description *string
	if desc, ok := request.GetArguments()["description"].(string); ok {
		description = &desc
	}
	
	var color *string
	if col, ok := request.GetArguments()["color"].(string); ok {
		color = &col
	}
	
	// Create category
	category, err := h.categoryService.CreateCategory(name, description, color)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}
	
	// Format response
	responseText := fmt.Sprintf("Category created successfully:\nID: %d\nName: %s", category.ID, category.Name)
	
	if category.Description != nil && *category.Description != "" {
		responseText += fmt.Sprintf("\nDescription: %s", *category.Description)
	}
	
	if category.Color != nil && *category.Color != "" {
		responseText += fmt.Sprintf("\nColor: %s", *category.Color)
	}
	
	responseText += fmt.Sprintf("\nCreated at: %s", category.CreatedAt.Format("2006-01-02 15:04:05"))
	
	return mcp.NewToolResultText(responseText), nil
}

// GetAllCategoriesHandler handles the get_all_categories MCP tool
func (h *CategoryHandler) GetAllCategoriesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	categories, err := h.categoryService.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve categories: %w", err)
	}
	
	if len(categories) == 0 {
		return mcp.NewToolResultText("No categories found"), nil
	}
	
	var responseText string
	for i, category := range categories {
		if i > 0 {
			responseText += "\n\n"
		}
		responseText += fmt.Sprintf("Category ID: %d\nName: %s", category.ID, category.Name)
		
		if category.Description != nil && *category.Description != "" {
			responseText += fmt.Sprintf("\nDescription: %s", *category.Description)
		}
		
		if category.Color != nil && *category.Color != "" {
			responseText += fmt.Sprintf("\nColor: %s", *category.Color)
		}
		
		responseText += fmt.Sprintf("\nCreated: %s", category.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	
	return mcp.NewToolResultText(responseText), nil
}

// GetCategoryHandler handles the get_category MCP tool
func (h *CategoryHandler) GetCategoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract category ID
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id parameter is required and must be a number")
	}
	id := int64(idRaw)
	
	category, err := h.categoryService.GetCategoryByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve category: %w", err)
	}
	
	responseText := fmt.Sprintf("Category Details:\nID: %d\nName: %s", category.ID, category.Name)
	
	if category.Description != nil && *category.Description != "" {
		responseText += fmt.Sprintf("\nDescription: %s", *category.Description)
	}
	
	if category.Color != nil && *category.Color != "" {
		responseText += fmt.Sprintf("\nColor: %s", *category.Color)
	}
	
	responseText += fmt.Sprintf("\nCreated: %s\nUpdated: %s", 
		category.CreatedAt.Format("2006-01-02 15:04:05"),
		category.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	return mcp.NewToolResultText(responseText), nil
}

// UpdateCategoryHandler handles the update_category MCP tool
func (h *CategoryHandler) UpdateCategoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract category ID
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id parameter is required and must be a number")
	}
	id := int64(idRaw)
	
	// Optional parameters
	var name *string
	if nameVal, ok := request.GetArguments()["name"].(string); ok {
		name = &nameVal
	}
	
	var description *string
	if desc, ok := request.GetArguments()["description"].(string); ok {
		description = &desc
	}
	
	var color *string
	if col, ok := request.GetArguments()["color"].(string); ok {
		color = &col
	}
	
	// Update category
	category, err := h.categoryService.UpdateCategory(id, name, description, color)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}
	
	responseText := fmt.Sprintf("Category updated successfully:\nID: %d\nName: %s", category.ID, category.Name)
	
	if category.Description != nil && *category.Description != "" {
		responseText += fmt.Sprintf("\nDescription: %s", *category.Description)
	}
	
	if category.Color != nil && *category.Color != "" {
		responseText += fmt.Sprintf("\nColor: %s", *category.Color)
	}
	
	responseText += fmt.Sprintf("\nUpdated at: %s", category.UpdatedAt.Format("2006-01-02 15:04:05"))
	
	return mcp.NewToolResultText(responseText), nil
}

// DeleteCategoryHandler handles the delete_category MCP tool
func (h *CategoryHandler) DeleteCategoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract category ID
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id parameter is required and must be a number")
	}
	id := int64(idRaw)
	
	err := h.categoryService.DeleteCategory(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %w", err)
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("Category with ID %d deleted successfully", id)), nil
}

// GetCategoryTodosHandler handles the get_category_todos MCP tool
func (h *CategoryHandler) GetCategoryTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract category ID
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("id parameter is required and must be a number")
	}
	id := int64(idRaw)
	
	todos, err := h.categoryService.GetTodosByCategory(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve todos for category: %w", err)
	}
	
	if len(todos) == 0 {
		return mcp.NewToolResultText("No todos found in this category"), nil
	}
	
	var resultText string
	for i, todo := range todos {
		if i > 0 {
			resultText += "\n\n"
		}
		resultText += fmt.Sprintf("Todo ID: %s\nTitle: %s", todo.ID, todo.Title)
		
		if todo.CompletedAt != nil {
			resultText += fmt.Sprintf("\nStatus: Completed (%s)", todo.CompletedAt.Format("2006-01-02 15:04:05"))
		} else {
			resultText += "\nStatus: Pending"
		}
		
		if todo.DueDate != nil {
			resultText += fmt.Sprintf("\nDue Date: %s", todo.DueDate.Format("2006-01-02 15:04:05"))
		}
	}
	
	return mcp.NewToolResultText(resultText), nil
}

// AssignTodoToCategoryHandler handles the assign_todo_to_category MCP tool
func (h *CategoryHandler) AssignTodoToCategoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	todoID, ok := request.GetArguments()["todo_id"].(string)
	if !ok || todoID == "" {
		return nil, fmt.Errorf("todo_id parameter is required and must be a string")
	}
	
	categoryIDRaw, ok := request.GetArguments()["category_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("category_id parameter is required and must be a number")
	}
	
	categoryID := int64(categoryIDRaw)
	
	// Call the todo service to assign the todo to the category
	todo, err := h.todoService.AssignTodoToCategory(todoID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to assign todo to category: %w", err)
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("Todo '%s' (ID: %s) successfully assigned to category %d", todo.Title, todo.ID, categoryID)), nil
}

// RemoveTodoFromCategoryHandler handles the remove_todo_from_category MCP tool
func (h *CategoryHandler) RemoveTodoFromCategoryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract parameters
	todoID, ok := request.GetArguments()["todo_id"].(string)
	if !ok || todoID == "" {
		return nil, fmt.Errorf("todo_id parameter is required and must be a string")
	}
	
	// Call the todo service to remove the todo from its category
	todo, err := h.todoService.RemoveTodoFromCategory(todoID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove todo from category: %w", err)
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("Todo '%s' (ID: %s) successfully removed from its category", todo.Title, todo.ID)), nil
}

// GetUncategorizedTodosHandler handles the get_uncategorized_todos MCP tool
func (h *CategoryHandler) GetUncategorizedTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todos, err := h.categoryService.GetUncategorizedTodos()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve uncategorized todos: %w", err)
	}
	
	if len(todos) == 0 {
		return mcp.NewToolResultText("No uncategorized todos found"), nil
	}
	
	var resultText string
	for i, todo := range todos {
		if i > 0 {
			resultText += "\n\n"
		}
		resultText += fmt.Sprintf("Todo ID: %s\nTitle: %s", todo.ID, todo.Title)
		
		if todo.CompletedAt != nil {
			resultText += fmt.Sprintf("\nStatus: Completed (%s)", todo.CompletedAt.Format("2006-01-02 15:04:05"))
		} else {
			resultText += "\nStatus: Pending"
		}
		
		if todo.DueDate != nil {
			resultText += fmt.Sprintf("\nDue Date: %s", todo.DueDate.Format("2006-01-02 15:04:05"))
		}
	}
	
	return mcp.NewToolResultText(resultText), nil
}