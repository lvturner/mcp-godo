package todo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAddTodo verifies adding a valid todo and error cases
func TestAddTodo(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		expectErr bool
	}{
		{"Valid title", "Buy groceries", false},
		{"Empty title", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service := &todo_in_memory{
				todos: make(map[string]TodoItem),
			}

			// Execute
			todo, err := service.AddTodo(tt.title)

			// Assert
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, service.todos, 1)
				assert.Equal(t, tt.title, service.todos[todo.ID].Title)
			}
		})
	}
}

// TestGetAllTodos verifies retrieving all todos
func TestGetAllTodos(t *testing.T) {
	service := &todo_in_memory{
		todos: map[string]TodoItem{
			"1": {ID: "1", Title: "Task 1"},
			"2": {ID: "2", Title: "Task 2"},
		},
	}

	todos := service.GetAllTodos()
	assert.Len(t, todos, 2)
	assert.Equal(t, "Task 1", todos[0].Title)
	assert.Equal(t, "Task 2", todos[1].Title)
}

// TestGetTodo verifies retrieving a todo by ID and invalid ID
func TestGetTodo(t *testing.T) {
	service := &todo_in_memory{
		todos: map[string]TodoItem{
			"1": {ID: "1", Title: "Task 1"},
		},
	}

	// Valid ID
	todo, err := service.GetTodo("1")
	assert.NoError(t, err)
	assert.Equal(t, "Task 1", todo.Title)

	// Invalid ID
	_, err = service.GetTodo("2")
	assert.Error(t, err)
}
