package todo

type TodoItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TodoService interface {
	AddTodo(title string) (TodoItem, error)
	GetAllTodos() []TodoItem
	GetTodo(id string) (TodoItem, error)
}
