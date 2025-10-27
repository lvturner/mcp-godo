package unit

import (
	"context"
	"testing"
	"time"

	"mcp-godo/pkg/handler"
	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService for testing
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) CreateProject(name string, description *string) (todo.Project, error) {
	args := m.Called(name, description)
	return args.Get(0).(todo.Project), args.Error(1)
}

func (m *MockProjectService) GetAllProjects() []todo.Project {
	args := m.Called()
	return args.Get(0).([]todo.Project)
}

func (m *MockProjectService) GetProject(id int64) (todo.Project, error) {
	args := m.Called(id)
	return args.Get(0).(todo.Project), args.Error(1)
}

func (m *MockProjectService) UpdateProject(id int64, name string, description *string) (todo.Project, error) {
	args := m.Called(id, name, description)
	return args.Get(0).(todo.Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(id int64) (todo.Project, error) {
	args := m.Called(id)
	return args.Get(0).(todo.Project), args.Error(1)
}

func (m *MockProjectService) GetProjectTodos(id int64) []todo.TodoItem {
	args := m.Called(id)
	return args.Get(0).([]todo.TodoItem)
}

// MockCategoryService for testing
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) CreateCategory(name string, description *string, color *string) (todo.Category, error) {
	args := m.Called(name, description, color)
	return args.Get(0).(todo.Category), args.Error(1)
}

func (m *MockCategoryService) GetAllCategories() ([]todo.Category, error) {
	args := m.Called()
	return args.Get(0).([]todo.Category), args.Error(1)
}

func (m *MockCategoryService) GetCategoryByID(id int64) (todo.Category, error) {
	args := m.Called(id)
	return args.Get(0).(todo.Category), args.Error(1)
}

func (m *MockCategoryService) UpdateCategory(id int64, name *string, description *string, color *string) (todo.Category, error) {
	args := m.Called(id, name, description, color)
	return args.Get(0).(todo.Category), args.Error(1)
}

func (m *MockCategoryService) DeleteCategory(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryService) GetTodosByCategory(categoryID int64) ([]todo.TodoItem, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]todo.TodoItem), args.Error(1)
}

func (m *MockCategoryService) GetUncategorizedTodos() ([]todo.TodoItem, error) {
	args := m.Called()
	return args.Get(0).([]todo.TodoItem), args.Error(1)
}

func TestGetActiveTodosHandler_NoTodos(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	// Mock no active todos
	mockTodoService.On("GetActiveTodos").Return([]todo.TodoItem{})
	
	ctx := context.Background()
	request := mcp.CallToolRequest{}
	
	result, err := h.GetActiveTodosHandler(ctx, request)
	
	assert.NoError(t, err)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "No active todos found")
	
	mockTodoService.AssertExpectations(t)
}

func TestGetActiveTodosHandler_WithTodos_NoProjectNoCategory(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	now := time.Now()
	todos := []todo.TodoItem{
		{
			ID:          "1",
			Title:       "Test Todo",
			CompletedAt: nil,
			DueDate:     &now,
			CreatedDate: now,
			ReferenceID: nil,
			ProjectID:   nil,
			CategoryID:  nil,
		},
	}
	
	mockTodoService.On("GetActiveTodos").Return(todos)
	
	ctx := context.Background()
	request := mcp.CallToolRequest{}
	
	result, err := h.GetActiveTodosHandler(ctx, request)
	
	assert.NoError(t, err)
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "ID: 1")
	assert.Contains(t, resultText, "Title: Test Todo")
	assert.Contains(t, resultText, "Status: Incomplete")
	assert.NotContains(t, resultText, "Project:")
	assert.NotContains(t, resultText, "Category:")
	
	mockTodoService.AssertExpectations(t)
}

func TestGetActiveTodosHandler_WithProjectAndCategory(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	now := time.Now()
	projectID := int64(1)
	categoryID := int64(2)
	
	todos := []todo.TodoItem{
		{
			ID:          "1",
			Title:       "Test Todo with Project and Category",
			CompletedAt: nil,
			DueDate:     &now,
			CreatedDate: now,
			ReferenceID: nil,
			ProjectID:   &projectID,
			CategoryID:  &categoryID,
		},
	}
	
	project := todo.Project{
		ID:   1,
		Name: "Test Project",
	}
	
	category := todo.Category{
		ID:   2,
		Name: "Test Category",
	}
	
	mockTodoService.On("GetActiveTodos").Return(todos)
	mockProjectService.On("GetProject", int64(1)).Return(project, nil)
	mockCategoryService.On("GetCategoryByID", int64(2)).Return(category, nil)
	
	ctx := context.Background()
	request := mcp.CallToolRequest{}
	
	result, err := h.GetActiveTodosHandler(ctx, request)
	
	assert.NoError(t, err)
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "ID: 1")
	assert.Contains(t, resultText, "Title: Test Todo with Project and Category")
	assert.Contains(t, resultText, "Status: Incomplete")
	assert.Contains(t, resultText, "Project: Test Project")
	assert.Contains(t, resultText, "Category: Test Category")
	
	mockTodoService.AssertExpectations(t)
	mockProjectService.AssertExpectations(t)
	mockCategoryService.AssertExpectations(t)
}

func TestGetActiveTodosHandler_WithProjectServiceError(t *testing.T) {
	mockTodoService := new(MockTodoService)
	mockProjectService := new(MockProjectService)
	mockCategoryService := new(MockCategoryService)
	
	h := handler.NewHandlerWithProjectAndCategory(mockTodoService, mockProjectService, mockCategoryService)
	
	now := time.Now()
	projectID := int64(1)
	categoryID := int64(2)
	
	todos := []todo.TodoItem{
		{
			ID:          "1",
			Title:       "Test Todo",
			CompletedAt: nil,
			DueDate:     &now,
			CreatedDate: now,
			ReferenceID: nil,
			ProjectID:   &projectID,
			CategoryID:  &categoryID,
		},
	}
	
	category := todo.Category{
		ID:   2,
		Name: "Test Category",
	}
	
	mockTodoService.On("GetActiveTodos").Return(todos)
	mockProjectService.On("GetProject", int64(1)).Return(todo.Project{}, assert.AnError)
	mockCategoryService.On("GetCategoryByID", int64(2)).Return(category, nil)
	
	ctx := context.Background()
	request := mcp.CallToolRequest{}
	
	result, err := h.GetActiveTodosHandler(ctx, request)
	
	assert.NoError(t, err)
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "ID: 1")
	assert.Contains(t, resultText, "Title: Test Todo")
	assert.NotContains(t, resultText, "Project:") // Should not show project due to error
	assert.Contains(t, resultText, "Category: Test Category")
	
	mockTodoService.AssertExpectations(t)
	mockProjectService.AssertExpectations(t)
	mockCategoryService.AssertExpectations(t)
}

func TestGetActiveTodosHandler_WithoutServices(t *testing.T) {
	mockTodoService := new(MockTodoService)
	
	// Create handler without project and category services
	h := handler.NewHandler(mockTodoService)
	
	now := time.Now()
	projectID := int64(1)
	categoryID := int64(2)
	
	todos := []todo.TodoItem{
		{
			ID:          "1",
			Title:       "Test Todo",
			CompletedAt: nil,
			DueDate:     &now,
			CreatedDate: now,
			ReferenceID: nil,
			ProjectID:   &projectID,
			CategoryID:  &categoryID,
		},
	}
	
	mockTodoService.On("GetActiveTodos").Return(todos)
	
	ctx := context.Background()
	request := mcp.CallToolRequest{}
	
	result, err := h.GetActiveTodosHandler(ctx, request)
	
	assert.NoError(t, err)
	resultText := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, resultText, "ID: 1")
	assert.Contains(t, resultText, "Title: Test Todo")
	assert.NotContains(t, resultText, "Project:") // Should not show project without service
	assert.NotContains(t, resultText, "Category:") // Should not show category without service
	
	mockTodoService.AssertExpectations(t)
}