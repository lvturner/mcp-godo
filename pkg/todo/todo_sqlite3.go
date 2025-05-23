package todo

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func NewTodoSQLite(db *sql.DB) TodoService {
	return &todo_sqlite{db: db}
}

type todo_sqlite struct {
	db *sql.DB
}

func (t *todo_sqlite) AddTodo(title string) (TodoItem, error) {
	if title == "" {
		return TodoItem{}, fmt.Errorf("title cannot be empty")
	}
	var id string
	err := t.db.QueryRow("INSERT INTO todos (title) VALUES ($1) RETURNING id", title).Scan(&id)
	if err != nil {
		return TodoItem{}, err
	}
	newItem := TodoItem{ID: id, Title: title}

	return newItem, nil
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

func (t *todo_sqlite) GetAllTodos() []TodoItem {
	rows, err := t.db.Query("SELECT id, title, completed FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var items []TodoItem
	for rows.Next() {
		var item TodoItem
		err = rows.Scan(&item.ID, &item.Title, &item.Completed)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, item)
	}
	return items
}

func (t *todo_sqlite) GetTodo(id string) (TodoItem, error) {
	var item TodoItem
	err := t.db.QueryRow("SELECT id, title, completed FROM todos WHERE id = $1", id).Scan(&item.ID, &item.Title, &item.Completed)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_sqlite) Close() error {
	return t.db.Close()
}