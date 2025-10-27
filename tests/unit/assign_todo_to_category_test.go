package unit

import (
	"context"
	"testing"

	"mcp-godo/pkg/handler"
	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// TestAssignTodoToCategoryHandler_Success tests successful assignment of todo to category
func TestAssignTodoToCategoryHandler_Success(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	expectedTodo := todo.TodoItem{
		ID:         "todo-123",
		Title:      "Test Todo",
		CategoryID: int64Ptr(42),
	}
	
	// Set up mock expectations
	mockTodoService.On("AssignTodoToCategory", "todo-123", int64(42)).Return(expectedTodo, nil)
	
	// Create request with valid parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "todo-123",
				"category_id": float64(42),
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify the response contains success message with todo title and category ID
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "Todo 'Test Todo' (ID: todo-123) successfully assigned to category 42")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_MissingTodoID tests missing todo_id parameter
func TestAssignTodoToCategoryHandler_MissingTodoID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	// Create request without todo_id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"category_id": float64(42),
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "todo_id parameter is required and must be a string")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_EmptyTodoID tests empty todo_id parameter
func TestAssignTodoToCategoryHandler_EmptyTodoID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	// Create request with empty todo_id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "",
				"category_id": float64(42),
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "todo_id parameter is required and must be a string")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_InvalidTodoIDType tests invalid todo_id type (not string)
func TestAssignTodoToCategoryHandler_InvalidTodoIDType(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	// Create request with todo_id as number instead of string
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     float64(123), // Invalid type
				"category_id": float64(42),
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "todo_id parameter is required and must be a string")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_MissingCategoryID tests missing category_id parameter
func TestAssignTodoToCategoryHandler_MissingCategoryID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	// Create request without category_id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id": "todo-123",
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "category_id parameter is required and must be a number")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_InvalidCategoryIDType tests invalid category_id type (not number)
func TestAssignTodoToCategoryHandler_InvalidCategoryIDType(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	// Create request with category_id as string instead of number
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "todo-123",
				"category_id": "42", // Invalid type
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "category_id parameter is required and must be a number")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_ZeroCategoryID tests zero category_id
func TestAssignTodoToCategoryHandler_ZeroCategoryID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	expectedTodo := todo.TodoItem{
		ID:         "todo-123",
		Title:      "Test Todo",
		CategoryID: int64Ptr(0),
	}
	
	// Set up mock expectations for zero category ID
	mockTodoService.On("AssignTodoToCategory", "todo-123", int64(0)).Return(expectedTodo, nil)
	
	// Create request with zero category_id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "todo-123",
				"category_id": float64(0),
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify the response contains success message with zero category ID
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "Todo 'Test Todo' (ID: todo-123) successfully assigned to category 0")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_NegativeCategoryID tests negative category_id
func TestAssignTodoToCategoryHandler_NegativeCategoryID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	expectedTodo := todo.TodoItem{
		ID:         "todo-123",
		Title:      "Test Todo",
		CategoryID: int64Ptr(-5),
	}
	
	// Set up mock expectations for negative category ID
	mockTodoService.On("AssignTodoToCategory", "todo-123", int64(-5)).Return(expectedTodo, nil)
	
	// Create request with negative category_id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "todo-123",
				"category_id": float64(-5),
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify the response contains success message with negative category ID
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "Todo 'Test Todo' (ID: todo-123) successfully assigned to category -5")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_ServiceError tests when todo service returns an error
func TestAssignTodoToCategoryHandler_ServiceError(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)
	
	// Set up mock to return error
	mockTodoService.On("AssignTodoToCategory", "todo-123", int64(42)).Return(todo.TodoItem{}, assert.AnError)
	
	// Create request with valid parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "todo-123",
				"category_id": float64(42),
			},
		},
	}
	
	ctx := context.Background()
	result, err := categoryHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to assign todo to category")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_Integration tests the full flow through main Handler
func TestAssignTodoToCategoryHandler_Integration(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockCategoryService := new(MockCategoryService)
	
	// Create main handler with category service
	mainHandler := handler.NewHandlerWithProjectAndCategory(mockTodoService, nil, mockCategoryService)
	
	expectedTodo := todo.TodoItem{
		ID:         "todo-123",
		Title:      "Test Todo",
		CategoryID: int64Ptr(42),
	}
	
	// Set up mock expectations
	mockTodoService.On("AssignTodoToCategory", "todo-123", int64(42)).Return(expectedTodo, nil)
	
	// Create request with valid parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "todo-123",
				"category_id": float64(42),
			},
		},
	}
	
	ctx := context.Background()
	result, err := mainHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify the response contains success message
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "Todo 'Test Todo' (ID: todo-123) successfully assigned to category 42")
	
	mockTodoService.AssertExpectations(t)
}

// TestAssignTodoToCategoryHandler_CategoryServiceNotInitialized tests when category service is nil
func TestAssignTodoToCategoryHandler_CategoryServiceNotInitialized(t *testing.T) {
	mockTodoService := new(MockTodoService)
	
	// Create main handler without category service
	mainHandler := handler.NewHandler(mockTodoService)
	
	// Create request with valid parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"todo_id":     "todo-123",
				"category_id": float64(42),
			},
		},
	}
	
	ctx := context.Background()
	result, err := mainHandler.AssignTodoToCategoryHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "category service not initialized")
	
	mockTodoService.AssertExpectations(t)
}