package unit

import (
	"context"
	"testing"

	"mcp-godo/pkg/handler"
	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProjectHandler_Success(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	expectedProject := todo.Project{
		ID:   1,
		Name: "Test Project",
	}
	
	// Set up mock expectations
	mockProjectService.On("DeleteProject", int64(1)).Return(expectedProject, nil)
	
	// Create request with valid parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": float64(1),
			},
		},
	}
	
	ctx := context.Background()
	result, err := h.DeleteProjectHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify the response contains success message with project details
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "Project deleted: ID=1, Name=Test Project")
	
	mockProjectService.AssertExpectations(t)
}

func TestDeleteProjectHandler_MissingID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	// Create request without id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}
	
	ctx := context.Background()
	result, err := h.DeleteProjectHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "project id is required and must be a number")
	
	mockProjectService.AssertExpectations(t)
}

func TestDeleteProjectHandler_InvalidIDType(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	// Create request with invalid id type (string instead of number)
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": "invalid",
			},
		},
	}
	
	ctx := context.Background()
	result, err := h.DeleteProjectHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "project id is required and must be a number")
	
	mockProjectService.AssertExpectations(t)
}

func TestDeleteProjectHandler_ProjectServiceNotInitialized(t *testing.T) {
	mockTodoService := new(MockTodoService)
	
	// Create handler without project service
	h := handler.NewHandler(mockTodoService)
	
	// Create request with valid parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": float64(1),
			},
		},
	}
	
	ctx := context.Background()
	result, err := h.DeleteProjectHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "project service not initialized")
	
	mockTodoService.AssertExpectations(t)
}

func TestDeleteProjectHandler_ServiceError(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	// Set up mock to return error
	mockProjectService.On("DeleteProject", int64(1)).Return(todo.Project{}, assert.AnError)
	
	// Create request with valid parameters
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": float64(1),
			},
		},
	}
	
	ctx := context.Background()
	result, err := h.DeleteProjectHandler(ctx, request)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to delete project")
	
	mockProjectService.AssertExpectations(t)
}

func TestDeleteProjectHandler_ZeroID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	expectedProject := todo.Project{
		ID:   0,
		Name: "Zero ID Project",
	}
	
	// Set up mock expectations for zero ID
	mockProjectService.On("DeleteProject", int64(0)).Return(expectedProject, nil)
	
	// Create request with zero id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": float64(0),
			},
		},
	}
	
	ctx := context.Background()
	result, err := h.DeleteProjectHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify the response contains success message with zero ID
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "Project deleted: ID=0, Name=Zero ID Project")
	
	mockProjectService.AssertExpectations(t)
}

func TestDeleteProjectHandler_NegativeID(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	expectedProject := todo.Project{
		ID:   -5,
		Name: "Negative ID Project",
	}
	
	// Set up mock expectations for negative ID
	mockProjectService.On("DeleteProject", int64(-5)).Return(expectedProject, nil)
	
	// Create request with negative id
	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": float64(-5),
			},
		},
	}
	
	ctx := context.Background()
	result, err := h.DeleteProjectHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	
	// Verify the response contains success message with negative ID
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "Project deleted: ID=-5, Name=Negative ID Project")
	
	mockProjectService.AssertExpectations(t)
}