package contract

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

// MockCategoryService for contract testing
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

// MockTodoService for contract testing
type MockTodoService struct {
	mock.Mock
}

func (m *MockTodoService) AddTodo(title string, dueDate *time.Time) (todo.TodoItem, error) {
	args := m.Called(title, dueDate)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) AddTodoToProject(title string, projectID int64, dueDate *time.Time) (todo.TodoItem, error) {
	args := m.Called(title, projectID, dueDate)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) AddTodoToCategory(title string, categoryID int64, dueDate *time.Time) (todo.TodoItem, error) {
	args := m.Called(title, categoryID, dueDate)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) GetAllTodos() []todo.TodoItem {
	args := m.Called()
	return args.Get(0).([]todo.TodoItem)
}

func (m *MockTodoService) GetActiveTodos() []todo.TodoItem {
	args := m.Called()
	return args.Get(0).([]todo.TodoItem)
}

func (m *MockTodoService) GetCompletedTodos() []todo.TodoItem {
	args := m.Called()
	return args.Get(0).([]todo.TodoItem)
}

func (m *MockTodoService) GetTodosByProject(projectID int64) []todo.TodoItem {
	args := m.Called(projectID)
	return args.Get(0).([]todo.TodoItem)
}

func (m *MockTodoService) GetTodosByCategory(categoryID int64) []todo.TodoItem {
	args := m.Called(categoryID)
	return args.Get(0).([]todo.TodoItem)
}

func (m *MockTodoService) GetUncategorizedTodos() []todo.TodoItem {
	args := m.Called()
	return args.Get(0).([]todo.TodoItem)
}

func (m *MockTodoService) GetTodo(id string) (todo.TodoItem, error) {
	args := m.Called(id)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) CompleteTodo(id string) (todo.TodoItem, error) {
	args := m.Called(id)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) UnCompleteTodo(id string) (todo.TodoItem, error) {
	args := m.Called(id)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) SetDueDate(id string, dueDate time.Time) (todo.TodoItem, error) {
	args := m.Called(id, dueDate)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) DeleteTodo(id string) (todo.TodoItem, error) {
	args := m.Called(id)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) TitleSearchTodo(query string, activeOnly bool) []todo.TodoItem {
	args := m.Called(query, activeOnly)
	return args.Get(0).([]todo.TodoItem)
}

func (m *MockTodoService) AssignTodoToCategory(todoID string, categoryID int64) (todo.TodoItem, error) {
	args := m.Called(todoID, categoryID)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) RemoveTodoFromCategory(todoID string) (todo.TodoItem, error) {
	args := m.Called(todoID)
	return args.Get(0).(todo.TodoItem), args.Error(1)
}

func (m *MockTodoService) AddRecurrencePattern(pattern todo.RecurrencePattern) (int64, error) {
	args := m.Called(pattern)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTodoService) GetRecurrencePatternByID(id int64) (todo.RecurrencePattern, error) {
	args := m.Called(id)
	return args.Get(0).(todo.RecurrencePattern), args.Error(1)
}

func TestCreateCategoryHandler_Success(t *testing.T) {
	mockCategoryService := new(MockCategoryService)
	mockTodoService := new(MockTodoService)
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)

	expectedCategory := todo.Category{
		ID:          1,
		Name:        "Work Tasks",
		Description: stringPtr("Professional tasks"),
		Color:       stringPtr("#3498db"),
	}

	mockCategoryService.On("CreateCategory", "Work Tasks", stringPtr("Professional tasks"), stringPtr("#3498db")).Return(expectedCategory, nil)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name":        "Work Tasks",
				"description": "Professional tasks",
				"color":       "#3498db",
			},
		},
	}

	result, err := categoryHandler.CreateCategoryHandler(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Category created successfully")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Work Tasks")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "#3498db")
	mockCategoryService.AssertExpectations(t)
}

func TestCreateCategoryHandler_MinimalParameters(t *testing.T) {
	mockCategoryService := new(MockCategoryService)
	mockTodoService := new(MockTodoService)
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)

	expectedCategory := todo.Category{
		ID:   1,
		Name: "Simple Category",
	}

	mockCategoryService.On("CreateCategory", "Simple Category", (*string)(nil), (*string)(nil)).Return(expectedCategory, nil)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name": "Simple Category",
			},
		},
	}

	result, err := categoryHandler.CreateCategoryHandler(context.Background(), request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Category created successfully")
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "Simple Category")
	mockCategoryService.AssertExpectations(t)
}

func TestCreateCategoryHandler_MissingName(t *testing.T) {
	mockCategoryService := new(MockCategoryService)
	mockTodoService := new(MockTodoService)
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"description": "This should fail",
			},
		},
	}

	result, err := categoryHandler.CreateCategoryHandler(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name parameter is required")
	mockCategoryService.AssertExpectations(t)
}

func TestCreateCategoryHandler_InvalidNameType(t *testing.T) {
	mockCategoryService := new(MockCategoryService)
	mockTodoService := new(MockTodoService)
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name": 123, // Invalid type
			},
		},
	}

	result, err := categoryHandler.CreateCategoryHandler(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name parameter is required and must be a string")
	mockCategoryService.AssertExpectations(t)
}

func TestCreateCategoryHandler_ServiceError(t *testing.T) {
	mockCategoryService := new(MockCategoryService)
	mockTodoService := new(MockTodoService)
	categoryHandler := handler.NewCategoryHandler(mockCategoryService, mockTodoService)

	mockCategoryService.On("CreateCategory", "Invalid Category", (*string)(nil), (*string)(nil)).Return(todo.Category{}, assert.AnError)

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"name": "Invalid Category",
			},
		},
	}

	result, err := categoryHandler.CreateCategoryHandler(context.Background(), request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to create category")
	mockCategoryService.AssertExpectations(t)
}

// Helper function
func stringPtr(s string) *string {
	return &s
}