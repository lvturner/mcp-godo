package todo

import (
	"time"
)

type TodoItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
	DueDate *time.Time `json:"due_date"`
	CreatedDate time.Time `json:"created_date"`
}

type TodoService interface {
	AddTodo(title string, dueDate *time.Time) (TodoItem, error)
	GetAllTodos() []TodoItem
	GetActiveTodos() []TodoItem
	GetCompletedTodos() []TodoItem
	GetTodo(id string) (TodoItem, error)
	CompleteTodo(id string) (TodoItem, error)
	UnCompleteTodo(id string) (TodoItem, error)
	SetDueDate(id string, dueDateStr time.Time) (TodoItem, error)
	DeleteTodo(id string) (TodoItem, error)
	TitleSearchTodo(query string, activeOnly bool) []TodoItem
}
