package contract

import (
	"context"
	"testing"

	"mcp-godo/pkg/handler"
	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

func TestGetAllCategoriesHandler_Success(t *testing.T) {
	mockService := new(MockCategoryService)
	categoryHandler := handler.NewCategoryHandler(mockService)

	expectedCategories := []todo.Category{
		{
			ID:          1,
			Name:        "Work Tasks",
			Description: stringPtr("Professional tasks"),
			Color:       stringPtr("#3498db"),
		},
		{
			ID:          2,
			Name:        "Personal Tasks",
			Description: stringPtr("Personal items"),
			Color:       stringPtr("#e74c3c"),
		},
	}

	mockService.On("GetAllCategories").Return(expectedCategories, nil)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := categoryHandler.GetAllCategoriesHandler(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Work Tasks")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Personal Tasks")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "#3498db")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "#e74c3c")
	mockService.AssertExpectations(t)
}

func TestGetAllCategoriesHandler_EmptyList(t *testing.T) {
	mockService := new(MockCategoryService)
	categoryHandler := handler.NewCategoryHandler(mockService)

	expectedCategories := []todo.Category{}

	mockService.On("GetAllCategories").Return(expectedCategories, nil)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := categoryHandler.GetAllCategoriesHandler(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "No categories found")
	mockService.AssertExpectations(t)
}

func TestGetCategoryHandler_Success(t *testing.T) {
	mockService := new(MockCategoryService)
	categoryHandler := handler.NewCategoryHandler(mockService)

	expectedCategory := todo.Category{
		ID:          1,
		Name:        "Work Tasks",
		Description: stringPtr("Professional tasks"),
		Color:       stringPtr("#3498db"),
	}

	mockService.On("GetCategoryByID", int64(1)).Return(expectedCategory, nil)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": 1.0,
			},
		},
	}

	result, err := categoryHandler.GetCategoryHandler(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Category Details")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Work Tasks")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "#3498db")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Professional tasks")
	mockService.AssertExpectations(t)
}

func TestGetCategoryHandler_MissingID(t *testing.T) {
	mockService := new(MockCategoryService)
	categoryHandler := handler.NewCategoryHandler(mockService)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{},
		},
	}

	result, err := categoryHandler.GetCategoryHandler(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "id parameter is required")
	mockService.AssertExpectations(t)
}

func TestGetCategoryHandler_InvalidIDType(t *testing.T) {
	mockService := new(MockCategoryService)
	categoryHandler := handler.NewCategoryHandler(mockService)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": "invalid", // String instead of number
			},
		},
	}

	result, err := categoryHandler.GetCategoryHandler(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "id parameter is required and must be a number")
	mockService.AssertExpectations(t)
}

func TestGetCategoryHandler_ServiceError(t *testing.T) {
	mockService := new(MockCategoryService)
	categoryHandler := handler.NewCategoryHandler(mockService)

	mockService.On("GetCategoryByID", int64(999)).Return(todo.Category{}, assert.AnError)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"id": 999.0,
			},
		},
	}

	result, err := categoryHandler.GetCategoryHandler(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to retrieve category")
	mockService.AssertExpectations(t)
}