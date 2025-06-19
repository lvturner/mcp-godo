package todo

import (
	"testing"
	"time"

	"fmt"

	"github.com/stretchr/testify/assert"
)

type mockTodoService struct {
	todos map[string]TodoItem
}

func newMockTodoService() *mockTodoService {
	return &mockTodoService{
		todos: make(map[string]TodoItem),
	}
}

func (m *mockTodoService) AddTodo(title string, dueDate *time.Time) (TodoItem, error) {
	id := "mock-" + title // Simple ID generation for testing
	item := TodoItem{
		ID:          id,
		Title:       title,
		Completed:   false,
		DueDate:     dueDate,
		CreatedDate: time.Now(),
	}
	m.todos[id] = item
	return item, nil
}

func (m *mockTodoService) GetAllTodos() []TodoItem {
	var items []TodoItem
	for _, item := range m.todos {
		items = append(items, item)
	}
	return items
}

func (m *mockTodoService) GetTodo(id string) (TodoItem, error) {
	item, exists := m.todos[id]
	if !exists {
		return TodoItem{}, fmt.Errorf("todo not found")
	}
	return item, nil
}

func TestTodoService(t *testing.T) {
	t.Run("AddTodo", func(t *testing.T) {
		svc := newMockTodoService()
		item, err := svc.AddTodo("test todo", nil)
		assert.NoError(t, err)
		assert.Equal(t, "test todo", item.Title)
		assert.False(t, item.Completed)
	})

	t.Run("TitleSearchTodo_Mock", func(t *testing.T) {
		svc := newMockTodoService()
		
		// Add test data
		testTitles := []string{
			"Buy groceries",
			"Finish project",
			"Buy new laptop",
			"Call mom",
			"Buy birthday gift",
		}
		for _, title := range testTitles {
			_, err := svc.AddTodo(title, nil)
			assert.NoError(t, err)
		}

		tests := []struct {
			name     string
			query    string
			expected []string
		}{
			{
				name:     "exact match",
				query:    "Buy groceries",
				expected: []string{"Buy groceries"},
			},
			{
				name:     "partial match",
				query:    "Buy",
				expected: []string{"Buy groceries", "Buy new laptop", "Buy birthday gift"},
			},
			{
				name:     "case insensitive",
				query:    "buy",
				expected: []string{"Buy groceries", "Buy new laptop", "Buy birthday gift"},
			},
			{
				name:     "no match",
				query:    "nonexistent",
				expected: []string{},
			},
			{
				name:     "empty query returns all",
				query:    "",
				expected: testTitles,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				results := svc.TitleSearchTodo(tt.query)
				assert.Equal(t, len(tt.expected), len(results))
				
				resultTitles := make([]string, len(results))
				for i, item := range results {
					resultTitles[i] = item.Title
				}
				
				for _, expected := range tt.expected {
					assert.Contains(t, resultTitles, expected)
				}
			})
		}
	})

	t.Run("GetTodo", func(t *testing.T) {
		svc := newMockTodoService()
		added, _ := svc.AddTodo("test", nil)
		
		got, err := svc.GetTodo(added.ID)
		assert.NoError(t, err)
		assert.Equal(t, added, got)
		
		_, err = svc.GetTodo("invalid")
		assert.Error(t, err)
	})

	t.Run("GetAllTodos", func(t *testing.T) {
		svc := newMockTodoService()
		svc.AddTodo("todo1", nil)
		svc.AddTodo("todo2", nil)

		all := svc.GetAllTodos()
		assert.Len(t, all, 2)
	})
}
