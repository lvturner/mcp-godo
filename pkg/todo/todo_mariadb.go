package todo

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewTodoMariaDB(db *sql.DB) TodoService {
	return &todo_mariadb{db: db}
}

type todo_mariadb struct {
	db *sql.DB
}

func (t *todo_mariadb) AddTodo(title string, dueDate *time.Time) (TodoItem, error) {
	if title == "" {
		return TodoItem{}, fmt.Errorf("title cannot be empty")
	}
	res, err := t.db.Exec("INSERT INTO todos (title, completed, due_date)  VALUES (?, ?, ?)", title, false, dueDate)
	if err != nil {
		return TodoItem{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return TodoItem{}, err
	}
	idStr := strconv.FormatInt(id, 10)
	newItem := TodoItem{ID: idStr, Title: title, Completed: false}
	if dueDate != nil {
		newItem.DueDate = dueDate
	}
	return newItem, nil
}

func (t *todo_mariadb) SetDueDate(id string, dueDate time.Time) (TodoItem, error) {
	_, err := t.db.Exec("UPDATE todos SET due_date = ? WHERE id = ?", dueDate, id)
	if err != nil {
		return TodoItem{}, err
	}
	item := TodoItem{ID: id}
	err = t.db.QueryRow("SELECT title, completed, due_date FROM todos WHERE id = ?", id).Scan(
		&item.Title, &item.Completed, &item.DueDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_mariadb) CompleteTodo(id string) (TodoItem, error) {
	_, err := t.db.Exec("UPDATE todos SET completed = true WHERE id = ?", id)
	if err != nil {
		return TodoItem{}, err
	}
	item := TodoItem{ID: id}
	err = t.db.QueryRow("SELECT title, completed, due_date FROM todos WHERE id = ?", id).Scan(
		&item.Title, &item.Completed, &item.DueDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_mariadb) UnCompleteTodo(id string) (TodoItem, error) {
	_, err := t.db.Exec("UPDATE todos SET completed = false WHERE id = ?", id)
	if err != nil {
		return TodoItem{}, err
	}
	item := TodoItem{ID: id}
	err = t.db.QueryRow("SELECT title, completed, due_date FROM todos WHERE id = ?", id).Scan(
		&item.Title, &item.Completed, &item.DueDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_mariadb) GetAllTodos() []TodoItem {
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

func (t *todo_mariadb) GetTodo(id string) (TodoItem, error) {
	var item TodoItem
	err := t.db.QueryRow("SELECT id, title, completed, due_date, created_date FROM todos WHERE id = ?", id).Scan(
		&item.ID, &item.Title, &item.Completed, &item.DueDate, &item.CreatedDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_mariadb) GetActiveTodos() []TodoItem {
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

func (t *todo_mariadb) GetCompletedTodos() []TodoItem {
	rows, err := t.db.Query("SELECT id, title, completed, due_date, created_date FROM todos WHERE completed = true")
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

func (t *todo_mariadb) DeleteTodo(id string) (TodoItem, error) {
	var item TodoItem
	row := t.db.QueryRow("SELECT id, title, completed, due_date FROM todos WHERE id = ?", id)
	err := row.Scan(&item.ID, &item.Title, &item.Completed, &item.DueDate)
	if err != nil {
		return item, err
	}
	_, err = t.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (t *todo_mariadb) Close() error {
	return t.db.Close()
}