package todo

type TodoItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Completed bool `json:"completed"`
}

type TodoService interface {
	AddTodo(title string) (TodoItem, error)
	GetAllTodos() []TodoItem
	GetTodo(id string) (TodoItem, error)
	CompleteTodo(id string) (TodoItem, error)
}
