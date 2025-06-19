package handler

import (
	"errors"
	"testing"
	"time"

	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

type mockTodoService struct {
	addTodoFunc            func(title string, dueDate *time.Time) (todo.TodoItem, error)
	getAllTodosFunc        func() []todo.TodoItem
	getActiveTodosFunc     func() []todo.TodoItem
	getCompletedTodosFunc  func() []todo.TodoItem
	getTodoFunc           func(id string) (todo.TodoItem, error)
	completeTodoFunc      func(id string) (todo.TodoItem, error)
	unCompleteTodoFunc    func(id string) (todo.TodoItem, error)
	setDueDateFunc        func(id string, dueDate time.Time) (todo.TodoItem, error)
	deleteTodoFunc        func(id string) (todo.TodoItem, error)
	titleSearchTodoFunc   func(query string) []todo.TodoItem
}

func (m *mockTodoService) AddTodo(title string, dueDate *time.Time) (todo.TodoItem, error) {
	return m.addTodoFunc(title, dueDate)
}

func (m *mockTodoService) GetAllTodos() []todo.TodoItem {
	return m.getAllTodosFunc()
}

func (m *mockTodoService) GetActiveTodos() []todo.TodoItem {
	return m.getActiveTodosFunc()
}

func (m *mockTodoService) GetCompletedTodos() []todo.TodoItem {
	return m.getCompletedTodosFunc()
}

func (m *mockTodoService) GetTodo(id string) (todo.TodoItem, error) {
	return m.getTodoFunc(id)
}

func (m *mockTodoService) CompleteTodo(id string) (todo.TodoItem, error) {
	return m.completeTodoFunc(id)
}

func (m *mockTodoService) UnCompleteTodo(id string) (todo.TodoItem, error) {
	return m.unCompleteTodoFunc(id)
}

func (m *mockTodoService) SetDueDate(id string, dueDate time.Time) (todo.TodoItem, error) {
	return m.setDueDateFunc(id, dueDate)
}

func (m *mockTodoService) DeleteTodo(id string) (todo.TodoItem, error) {
	return m.deleteTodoFunc(id)
}

func (m *mockTodoService) TitleSearchTodo(query string) []todo.TodoItem {
	return m.titleSearchTodoFunc(query)
}

func (m *mockTodoService) Close() error {
	return nil
}

func TestAddTodoHandler(t *testing.T) {
	tests := []struct {
		name          string
		args         map[string]interface{}
		mockFunc     func(title string, dueDate *time.Time) (todo.TodoItem, error)
		expectedText string
		expectError  bool
	}{
		{
			name: "success without due date",
			args: map[string]interface{}{"title": "test"},
			mockFunc: func(title string, dueDate *time.Time) (todo.TodoItem, error) {
				return todo.TodoItem{Title: title}, nil
			},
			expectedText: "test added to todo list",
			expectError: false,
		},
		{
			name: "success with due date",
			args: map[string]interface{}{
				"title": "test",
				"due_date": "2023-01-01T00:00:00Z",
			},
			mockFunc: func(title string, dueDate *time.Time) (todo.TodoItem, error) {
				return todo.TodoItem{Title: title}, nil
			},
			expectedText: "test added to todo list",
			expectError: false,
		},
		{
			name: "invalid title",
			args: map[string]interface{}{"title": 123},
			mockFunc: func(title string, dueDate *time.Time) (todo.TodoItem, error) {
				return todo.TodoItem{}, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockTodoService{
				addTodoFunc: tt.mockFunc,
			}
			h := NewHandler(mockSvc)
			
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := h.AddTodoHandler(nil, req)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedText, result.Content[0].(mcp.TextContent).Text)
			}
		})
	}
}

func TestCompleteTodoHandler(t *testing.T) {
	tests := []struct {
		name          string
		args         map[string]interface{}
		mockFunc     func(id string) (todo.TodoItem, error)
		expectedText string
		expectError  bool
	}{
		{
			name: "success",
			args: map[string]interface{}{"id": "123"},
			mockFunc: func(id string) (todo.TodoItem, error) {
				return todo.TodoItem{ID: id, Title: "test"}, nil
			},
			expectedText: "Todo test completed",
			expectError: false,
		},
		{
			name: "invalid id",
			args: map[string]interface{}{"id": 123},
			mockFunc: func(id string) (todo.TodoItem, error) {
				return todo.TodoItem{}, nil
			},
			expectError: true,
		},
		{
			name: "service error",
			args: map[string]interface{}{"id": "123"},
			mockFunc: func(id string) (todo.TodoItem, error) {
				return todo.TodoItem{}, errors.New("service error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockTodoService{
				completeTodoFunc: tt.mockFunc,
			}
			h := NewHandler(mockSvc)
			
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := h.CompleteTodoHandler(nil, req)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedText, result.Content[0].(mcp.TextContent).Text)
			}
		})
	}
}

func TestListTodosHandler(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		mockFunc     func() []todo.TodoItem
		expectError  bool
	}{
		{
			name: "success with todos",
			mockFunc: func() []todo.TodoItem {
				return []todo.TodoItem{
					{ID: "1", Title: "test1", Completed: false, CreatedDate: now},
					{ID: "2", Title: "test2", Completed: true, CreatedDate: now},
				}
			},
			expectError: false,
		},
		{
			name: "success empty",
			mockFunc: func() []todo.TodoItem {
				return []todo.TodoItem{}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockTodoService{
				getAllTodosFunc: tt.mockFunc,
			}
			h := NewHandler(mockSvc)
			
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{},
			}

			result, err := h.ListTodosHandler(nil, req)
			
			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

// Similar test structures can be created for other handler methods:
// - GetTodoHandler
// - DeleteTodoHandler  
// - GetActiveTodosHandler
// - GetCompletedTodosHandler
// - UnCompleteTodoHandler
// - UpdateDueDateHandler
// - TitleSearchHandler
