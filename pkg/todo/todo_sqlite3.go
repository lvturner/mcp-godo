package todo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewTodoSQLite(db *sql.DB) TodoService {
	return &todo_sqlite{db: db}
}

type todo_sqlite struct {
	db *sql.DB
}

func (t *todo_sqlite) AddTodo(title string, dueDate *time.Time) (TodoItem, error) {
	if title == "" {
		return TodoItem{}, fmt.Errorf("title cannot be empty")
	}
	var id string
	err := t.db.QueryRow("INSERT INTO todos (title, completed, due_date) VALUES ($1, false, $2) RETURNING id", title).Scan(&id)
	if err != nil {
		return TodoItem{}, err
	}
	newItem := TodoItem{ID: id, Title: title}
	if dueDate != nil {
		newItem.DueDate = dueDate
	}

	return newItem, nil
}

func (t *todo_sqlite) SetDueDate(id string, dueDate time.Time) (TodoItem, error) {
	_, err := t.db.Exec("UPDATE todos SET due_date = $1 WHERE id = $2", dueDate, id)
	if err != nil {
		return TodoItem{}, err
	}
	item := TodoItem{ID: id}
	err = t.db.QueryRow("SELECT title, completed, due_date, created_date FROM todos WHERE id = $1", id).Scan(&item.Title, &item.Completed, &item.DueDate, &item.CreatedDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_sqlite) CompleteTodo(id string) (TodoItem, error) {
	_, err := t.db.Exec("UPDATE todos SET completed = true WHERE id = $1", id)
	if err != nil {
		return TodoItem{}, err
	}
	item := TodoItem{ID: id}
	err = t.db.QueryRow("SELECT title, completed FROM todos WHERE id = $1", id).Scan(&item.Title, &item.Completed)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t* todo_sqlite) UnCompleteTodo(id string) (TodoItem, error) {
	_, err := t.db.Exec("UPDATE todos SET completed = false WHERE id = $1", id)
	if err != nil {
		return TodoItem{}, err
	}
	item := TodoItem{ID: id}
	err = t.db.QueryRow("SELECT title, completed FROM todos WHERE id = $1", id).Scan(&item.Title, &item.Completed)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_sqlite) GetAllTodos() []TodoItem {
	rows, err := t.db.Query("SELECT id, title, completed, due_date, created_date FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var items []TodoItem
	for rows.Next() {
		var item TodoItem
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &item.DueDate, &item.CreatedDate)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items
}

func (t *todo_sqlite) GetTodo(id string) (TodoItem, error) {
	var item TodoItem
	err := t.db.QueryRow("SELECT id, title, completed, due_date, created_date FROM todos WHERE id = $1", id).Scan(&item.ID, &item.Title, &item.Completed, &item.DueDate, &item.CreatedDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_sqlite) GetActiveTodos() []TodoItem {
	rows, err := t.db.Query("SELECT id, title, completed, due_date, created_date FROM todos WHERE completed = false")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var items []TodoItem
	for rows.Next() {
		var item TodoItem
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &item.DueDate, &item.CreatedDate)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items
}

func (t *todo_sqlite) GetCompletedTodos() []TodoItem {
	rows, err := t.db.Query("SELECT id, title, completed, created_date FROM todos WHERE completed = true")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var items []TodoItem
	for rows.Next() {
		var item TodoItem
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &item.CreatedDate)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items
}

func (t *todo_sqlite) DeleteTodo(id string) (TodoItem, error) {
	var item TodoItem
	row := t.db.QueryRow("SELECT id, title, completed, created_date FROM todos WHERE id = ?", id)
	err := row.Scan(&item.ID, &item.Title, &item.Completed, &item.CreatedDate)
	if err != nil {
		return item, err
	}
	_, err = t.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (t *todo_sqlite) TitleSearchTodo(query string) []TodoItem {
	rows, err := t.db.Query("SELECT id, title, completed, created_date FROM todos WHERE title LIKE ?", "%"+query+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var items []TodoItem
	for rows.Next() {
		var item TodoItem
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &item.CreatedDate)
		if err != nil {				
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items
}

func (t *todo_sqlite) Close() error {
	return t.db.Close()
}
