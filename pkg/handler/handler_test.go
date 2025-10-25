package handler

import (
	"errors"
	"fmt"
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
	titleSearchTodoFunc   func(query string, activeOnly bool) []todo.TodoItem
	addRecurrencePatternFunc    func(pattern todo.RecurrencePattern) (int64, error)
	getRecurrencePatternByIDFunc func(id int64) (todo.RecurrencePattern, error)
	addTodoToProjectFunc  func(title string, projectID int64, dueDate *time.Time) (todo.TodoItem, error)
	addTodoToCategoryFunc func(title string, categoryID int64, dueDate *time.Time) (todo.TodoItem, error)
	getTodosByProjectFunc func(projectID int64) []todo.TodoItem
	getTodosByCategoryFunc func(categoryID int64) []todo.TodoItem
	getUncategorizedTodosFunc func() []todo.TodoItem
	assignTodoToCategoryFunc func(todoID string, categoryID int64) (todo.TodoItem, error)
	removeTodoFromCategoryFunc func(todoID string) (todo.TodoItem, error)
}

func TestAddRecurrencePatternHandler(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name          string
		args         map[string]interface{}
		mockFunc     func(pattern todo.RecurrencePattern) (int64, error)
		expectedText string
		expectError  bool
	}{
		{
			name: "success with until and count",
			args: map[string]interface{}{
				"todo_id":   "123",
				"frequency": "weekly",
				"interval":  1.0,
				"until":     now.Format(time.RFC3339),
				"count":     5.0,
			},
			mockFunc: func(pattern todo.RecurrencePattern) (int64, error) {
				if pattern.TodoID != "123" || pattern.Frequency != "weekly" || pattern.Interval != 1 || pattern.Until == nil || pattern.Count == nil || *pattern.Count != 5 {
					return 0, fmt.Errorf("unexpected pattern")
				}
				return 1, nil
			},
			expectedText: "Recurrence pattern added with ID: 1",
			expectError: false,
		},
		{
			name: "success without until and count",
			args: map[string]interface{}{
				"todo_id":   "123",
				"frequency": "daily",
				"interval":  2.0,
			},
			mockFunc: func(pattern todo.RecurrencePattern) (int64, error) {
				if pattern.TodoID != "123" || pattern.Frequency != "daily" || pattern.Interval != 2 || pattern.Until != nil || pattern.Count != nil {
					return 0, fmt.Errorf("unexpected pattern")
				}
				return 2, nil
			},
			expectedText: "Recurrence pattern added with ID: 2",
			expectError: false,
		},
		{
			name: "invalid todo_id",
			args: map[string]interface{}{
				"frequency": "weekly",
				"interval":  1.0,
			},
			mockFunc: func(pattern todo.RecurrencePattern) (int64, error) {
				return 0, nil
			},
			expectError: true,
		},
		{
			name: "invalid frequency",
			args: map[string]interface{}{
				"todo_id":   "123",
				"interval":  1.0,
			},
			mockFunc: func(pattern todo.RecurrencePattern) (int64, error) {
				return 0, nil
			},
			expectError: true,
		},
		{
			name: "invalid interval",
			args: map[string]interface{}{
				"todo_id":   "123",
				"frequency": "weekly",
			},
			mockFunc: func(pattern todo.RecurrencePattern) (int64, error) {
				return 0, nil
			},
			expectError: true,
		},
		{
			name: "service error",
			args: map[string]interface{}{
				"todo_id":   "123",
				"frequency": "weekly",
				"interval":  1.0,
			},
			mockFunc: func(pattern todo.RecurrencePattern) (int64, error) {
				return 0, fmt.Errorf("service error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockTodoService{
				addRecurrencePatternFunc: tt.mockFunc,
			}
			h := NewHandler(mockSvc)
			
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := h.AddRecurrencePatternHandler(nil, req)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedText, result.Content[0].(mcp.TextContent).Text)
			}
		})
	}
}

func TestGetRecurrencePatternHandler(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name          string
		args         map[string]interface{}
		mockFunc     func(id int64) (todo.RecurrencePattern, error)
		expectedText string
		expectError  bool
	}{
		{
			name: "success with until and count",
			args: map[string]interface{}{"id": 1.0},
			mockFunc: func(id int64) (todo.RecurrencePattern, error) {
				count := 5
				return todo.RecurrencePattern{
					ID:        1,
					TodoID:    "123",
					Frequency: "weekly",
					Interval:  1,
					Until:     &now,
					Count:     &count,
				}, nil
			},
			expectedText: fmt.Sprintf("ID: 1, TodoID: 123, Frequency: weekly, Interval: 1, Until: %s, Count: 5", now.Format(time.RFC3339)),
			expectError: false,
		},
		{
			name: "success without until and count",
			args: map[string]interface{}{"id": 2.0},
			mockFunc: func(id int64) (todo.RecurrencePattern, error) {
				return todo.RecurrencePattern{
					ID:        2,
					TodoID:    "456",
					Frequency: "daily",
					Interval:  2,
				}, nil
			},
			expectedText: "ID: 2, TodoID: 456, Frequency: daily, Interval: 2",
			expectError: false,
		},
		{
			name: "invalid id",
			args: map[string]interface{}{},
			mockFunc: func(id int64) (todo.RecurrencePattern, error) {
				return todo.RecurrencePattern{}, nil
			},
			expectError: true,
		},
		{
			name: "service error",
			args: map[string]interface{}{"id": 3.0},
			mockFunc: func(id int64) (todo.RecurrencePattern, error) {
				return todo.RecurrencePattern{}, fmt.Errorf("service error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockTodoService{
				getRecurrencePatternByIDFunc: tt.mockFunc,
			}
			h := NewHandler(mockSvc)
			
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := h.GetRecurrencePatternHandler(nil, req)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedText, result.Content[0].(mcp.TextContent).Text)
			}
		})
	}
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

func (m *mockTodoService) TitleSearchTodo(query string, activeOnly bool) []todo.TodoItem {
	return m.titleSearchTodoFunc(query, activeOnly)
}

func (m *mockTodoService) AddRecurrencePattern(pattern todo.RecurrencePattern) (int64, error) {
	return m.addRecurrencePatternFunc(pattern)
}

func (m *mockTodoService) GetRecurrencePatternByID(id int64) (todo.RecurrencePattern, error) {
	return m.getRecurrencePatternByIDFunc(id)
}

func (m *mockTodoService) AddTodoToProject(title string, projectID int64, dueDate *time.Time) (todo.TodoItem, error) {
	if m.addTodoToProjectFunc != nil {
		return m.addTodoToProjectFunc(title, projectID, dueDate)
	}
	return todo.TodoItem{}, nil
}

func (m *mockTodoService) AddTodoToCategory(title string, categoryID int64, dueDate *time.Time) (todo.TodoItem, error) {
	if m.addTodoToCategoryFunc != nil {
		return m.addTodoToCategoryFunc(title, categoryID, dueDate)
	}
	return todo.TodoItem{}, nil
}

func (m *mockTodoService) GetTodosByProject(projectID int64) []todo.TodoItem {
	if m.getTodosByProjectFunc != nil {
		return m.getTodosByProjectFunc(projectID)
	}
	return []todo.TodoItem{}
}

func (m *mockTodoService) GetTodosByCategory(categoryID int64) []todo.TodoItem {
	if m.getTodosByCategoryFunc != nil {
		return m.getTodosByCategoryFunc(categoryID)
	}
	return []todo.TodoItem{}
}

func (m *mockTodoService) GetUncategorizedTodos() []todo.TodoItem {
	if m.getUncategorizedTodosFunc != nil {
		return m.getUncategorizedTodosFunc()
	}
	return []todo.TodoItem{}
}

func (m *mockTodoService) AssignTodoToCategory(todoID string, categoryID int64) (todo.TodoItem, error) {
	if m.assignTodoToCategoryFunc != nil {
		return m.assignTodoToCategoryFunc(todoID, categoryID)
	}
	return todo.TodoItem{}, nil
}

func (m *mockTodoService) RemoveTodoFromCategory(todoID string) (todo.TodoItem, error) {
	if m.removeTodoFromCategoryFunc != nil {
		return m.removeTodoFromCategoryFunc(todoID)
	}
	return todo.TodoItem{}, nil
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
					{ID: "1", Title: "test1", CompletedAt: nil, CreatedDate: now},
					{ID: "2", Title: "test2", CompletedAt: &now, CreatedDate: now},
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

func TestTitleSearchHandler(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		args         map[string]interface{}
		mockFunc     func(query string, activeOnly bool) []todo.TodoItem
		expectedText string
		expectError  bool
	}{
		{
			name: "success with results",
			args: map[string]interface{}{"query": "test"},
			mockFunc: func(query string, activeOnly bool) []todo.TodoItem {
				return []todo.TodoItem{
					{ID: "1", Title: "test todo", CompletedAt: nil, CreatedDate: now},
					{ID: "2", Title: "another test", CompletedAt: &now, CreatedDate: now},
				}
			},
			expectedText: fmt.Sprintf("ID: 1, Title: test todo, CompletedAt: <nil>, Due Date: \nID: 2, Title: another test, CompletedAt: %s, Due Date: ", now),
			expectError: false,
		},
		{
			name: "success with activeOnly filter",
			args: map[string]interface{}{"query": "test", "active_only": true},
			mockFunc: func(query string, activeOnly bool) []todo.TodoItem {
				if activeOnly {
					return []todo.TodoItem{
						{ID: "1", Title: "test todo", CompletedAt: nil, CreatedDate: now},
					}
				}
				return []todo.TodoItem{
					{ID: "1", Title: "test todo", CompletedAt: nil, CreatedDate: now},
					{ID: "2", Title: "another test", CompletedAt: &now, CreatedDate: now},
				}
			},
			expectedText: "ID: 1, Title: test todo, CompletedAt: <nil>, Due Date: ",
			expectError: false,
		},
		{
			name: "success no results",
			args: map[string]interface{}{"query": "nonexistent"},
			mockFunc: func(query string, activeOnly bool) []todo.TodoItem {
				return []todo.TodoItem{}
			},
			expectedText: "",
			expectError: false,
		},
		{
			name: "missing query param",
			args: map[string]interface{}{},
			mockFunc: func(query string, activeOnly bool) []todo.TodoItem {
				return []todo.TodoItem{}
			},
			expectError: true,
		},
		{
			name: "invalid query type",
			args: map[string]interface{}{"query": 123},
			mockFunc: func(query string, activeOnly bool) []todo.TodoItem {
				return []todo.TodoItem{}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockTodoService{
				titleSearchTodoFunc: tt.mockFunc,
			}
			h := NewHandler(mockSvc)
			
			req := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Arguments: tt.args,
				},
			}

			result, err := h.TitleSearchHandler(nil, req)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedText == "" {
					assert.Equal(t, "No todos found", result.Content[0].(mcp.TextContent).Text)
				} else {
					assert.Equal(t, tt.expectedText, result.Content[0].(mcp.TextContent).Text)
				}
			}
		})
	}
}
