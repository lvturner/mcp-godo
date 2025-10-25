package todo

import (
	"time"
)

type RecurrencePattern struct {
	ID        int64      `json:"id"`
	TodoID    string     `json:"todo_id"`
	Frequency string     `json:"frequency"`
	Interval  int        `json:"interval"`
	Until     *time.Time `json:"until"`
	Count     *int       `json:"count"`
}

type TodoItem struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completed_at"` // nil means not completed
	DueDate     *time.Time `json:"due_date"`
	CreatedDate time.Time  `json:"created_date"`
	ReferenceID *int64     `json:"reference_id"` // pointer to handle NULL in database
	ProjectID   *int64     `json:"project_id"`   // pointer to handle NULL in database (optional project association)
}

type TodoService interface {
	AddTodo(title string, dueDate *time.Time) (TodoItem, error)
	AddTodoToProject(title string, projectID int64, dueDate *time.Time) (TodoItem, error)
	GetAllTodos() []TodoItem
	GetActiveTodos() []TodoItem
	GetCompletedTodos() []TodoItem
	GetTodosByProject(projectID int64) []TodoItem
	GetTodo(id string) (TodoItem, error)
	CompleteTodo(id string) (TodoItem, error)
	UnCompleteTodo(id string) (TodoItem, error)
	SetDueDate(id string, dueDateStr time.Time) (TodoItem, error)
	DeleteTodo(id string) (TodoItem, error)
	TitleSearchTodo(query string, activeOnly bool) []TodoItem

	AddRecurrencePattern(pattern RecurrencePattern) (int64, error)
	GetRecurrencePatternByID(id int64) (RecurrencePattern, error)
}
