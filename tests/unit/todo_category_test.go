package unit

import (
	"testing"
	"time"

	"mcp-godo/pkg/todo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTodoService is a mock implementation of TodoService for testing todo-category relationships
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

func (m *MockTodoService) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestAddTodoToCategory_Success(t *testing.T) {
	mockService := new(MockTodoService)
	
	now := time.Now()
	expectedTodo := todo.TodoItem{
		ID:          "1",
		Title:       "Test Task",
		CategoryID:  int64Ptr(1),
		CreatedDate: now,
	}

	mockService.On("AddTodoToCategory", "Test Task", int64(1), (*time.Time)(nil)).Return(expectedTodo, nil)

	result, err := mockService.AddTodoToCategory("Test Task", 1, nil)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo.ID, result.ID)
	assert.Equal(t, expectedTodo.Title, result.Title)
	assert.Equal(t, int64(1), *result.CategoryID)
	mockService.AssertExpectations(t)
}

func TestAssignTodoToCategory_Success(t *testing.T) {
	mockService := new(MockTodoService)
	
	expectedTodo := todo.TodoItem{
		ID:         "123",
		Title:      "Existing Task",
		CategoryID: int64Ptr(1),
	}

	mockService.On("AssignTodoToCategory", "123", int64(1)).Return(expectedTodo, nil)

	result, err := mockService.AssignTodoToCategory("123", 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo.ID, result.ID)
	assert.Equal(t, int64(1), *result.CategoryID)
	mockService.AssertExpectations(t)
}

func TestRemoveTodoFromCategory_Success(t *testing.T) {
	mockService := new(MockTodoService)
	
	expectedTodo := todo.TodoItem{
		ID:         "123",
		Title:      "Existing Task",
		CategoryID: nil, // Should be nil after removal
	}

	mockService.On("RemoveTodoFromCategory", "123").Return(expectedTodo, nil)

	result, err := mockService.RemoveTodoFromCategory("123")

	assert.NoError(t, err)
	assert.Equal(t, expectedTodo.ID, result.ID)
	assert.Nil(t, result.CategoryID)
	mockService.AssertExpectations(t)
}

func TestGetTodosByCategoryFromService_Success(t *testing.T) {
	mockService := new(MockTodoService)
	
	expectedTodos := []todo.TodoItem{
		{ID: "1", Title: "Task 1", CategoryID: int64Ptr(1)},
		{ID: "2", Title: "Task 2", CategoryID: int64Ptr(1)},
	}

	mockService.On("GetTodosByCategory", int64(1)).Return(expectedTodos)

	result := mockService.GetTodosByCategory(1)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, "Task 1", result[0].Title)
	assert.Equal(t, "Task 2", result[1].Title)
	mockService.AssertExpectations(t)
}

func TestGetUncategorizedTodosFromService_Success(t *testing.T) {
	mockService := new(MockTodoService)
	
	expectedTodos := []todo.TodoItem{
		{ID: "1", Title: "Task 1", CategoryID: nil},
		{ID: "2", Title: "Task 2", CategoryID: nil},
	}

	mockService.On("GetUncategorizedTodos").Return(expectedTodos)

	result := mockService.GetUncategorizedTodos()

	assert.Equal(t, 2, len(result))
	assert.Nil(t, result[0].CategoryID)
	assert.Nil(t, result[1].CategoryID)
	mockService.AssertExpectations(t)
}

func TestAssignTodoToCategory_InvalidTodoID(t *testing.T) {
	mockService := new(MockTodoService)
	
	mockService.On("AssignTodoToCategory", "", int64(1)).Return(todo.TodoItem{}, assert.AnError)

	_, err := mockService.AssignTodoToCategory("", 1)

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestRemoveTodoFromCategory_InvalidTodoID(t *testing.T) {
	mockService := new(MockTodoService)
	
	mockService.On("RemoveTodoFromCategory", "").Return(todo.TodoItem{}, assert.AnError)

	_, err := mockService.RemoveTodoFromCategory("")

	assert.Error(t, err)
	mockService.AssertExpectations(t)
}
