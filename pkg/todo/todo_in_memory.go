package todo

import (
	"fmt"

	"github.com/google/uuid"
)

type todo_in_memory struct {
	todos map[string]TodoItem
}

func NewTodoInMemory() TodoService {
	return &todo_in_memory{todos: make(map[string]TodoItem)}
}

func (t *todo_in_memory) AddTodo(title string) (TodoItem, error) {
	if title == "" {
		return TodoItem{}, fmt.Errorf("title cannot be empty")
	}

	id := uuid.New().String()
	newItem := TodoItem{ID: id, Title: title}
	t.todos[id] = newItem

	return newItem, nil
}

func (t *todo_in_memory) CompleteTodo(id string) (TodoItem, error) {
	if _, exists := t.todos[id]; !exists {
		return TodoItem{}, fmt.Errorf("todo not found")
	}

	item := t.todos[id]
	item.Completed = true
	t.todos[id] = item

	return t.todos[id], nil
}

func (t *todo_in_memory) GetAllTodos() []TodoItem {
	items := make([]TodoItem, 0, len(t.todos))
	for _, item := range t.todos {
		items = append(items, item)
	}
	return items
}

func (t *todo_in_memory) GetTodo(id string) (TodoItem, error) {
	if todo, exists := t.todos[id]; exists {
		return todo, nil
	}
	return TodoItem{}, fmt.Errorf("todo not found")
}
