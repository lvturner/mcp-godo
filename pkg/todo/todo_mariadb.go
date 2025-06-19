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
	stmt, err := t.db.Prepare("INSERT INTO todos (title, completed, due_date) VALUES (?, ?, ?)")
	if err != nil {
		return TodoItem{}, err
	}
	defer stmt.Close()
	
	res, err := stmt.Exec(title, false, dueDate)
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
	stmt, err := t.db.Prepare("UPDATE todos SET due_date = ? WHERE id = ?")
	if err != nil {
		return TodoItem{}, err
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(dueDate, id)
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
	stmt, err := t.db.Prepare("UPDATE todos SET completed = true WHERE id = ?")
	if err != nil {
		return TodoItem{}, err
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(id)
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
	stmt, err := t.db.Prepare("UPDATE todos SET completed = false WHERE id = ?")
	if err != nil {
		return TodoItem{}, err
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(id)
	if err != nil {
		return TodoItem{}, err
	}
	var item TodoItem
	err = t.db.QueryRow("SELECT id, title, completed, due_date FROM todos WHERE id = ?", id).Scan(
		&item.ID, &item.Title, &item.Completed, &item.DueDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_mariadb) GetAllTodos() []TodoItem {
	stmt, err := t.db.Prepare("SELECT id, title, completed, due_date, created_date FROM todos")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
	rows, err := stmt.Query()
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
	stmt, err := t.db.Prepare("SELECT id, title, completed, due_date, created_date FROM todos WHERE id = ?")
	if err != nil {
		return TodoItem{}, err
	}
	defer stmt.Close()
	
	err = stmt.QueryRow(id).Scan(
		&item.ID, &item.Title, &item.Completed, &item.DueDate, &item.CreatedDate)
	if err != nil {
		return TodoItem{}, err
	}
	return item, nil
}

func (t *todo_mariadb) GetActiveTodos() []TodoItem {
	stmt, err := t.db.Prepare("SELECT id, title, completed, due_date, created_date FROM todos WHERE completed = false")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
	rows, err := stmt.Query()
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
	stmt, err := t.db.Prepare("SELECT id, title, completed, due_date, created_date FROM todos WHERE completed = true")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
	rows, err := stmt.Query()
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
	stmt, err := t.db.Prepare("SELECT id, title, completed, due_date FROM todos WHERE id = ?")
	if err != nil {
		return item, err
	}
	defer stmt.Close()
	
	row := stmt.QueryRow(id)
	err = row.Scan(&item.ID, &item.Title, &item.Completed, &item.DueDate)
	if err != nil {
		return item, err
	}
	stmt, err = t.db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		return item, err
	}
	defer stmt.Close()
	
	_, err = stmt.Exec(id)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (t *todo_mariadb) TitleSearchTodo(query string, activeOnly bool) []TodoItem {
	var queryStr string
	if activeOnly {
		queryStr = "SELECT id, title, completed, due_date, created_date FROM todos WHERE title LIKE ? AND completed = false"
	} else {
		queryStr = "SELECT id, title, completed, due_date, created_date FROM todos WHERE title LIKE ?"
	}

	stmt, err := t.db.Prepare(queryStr)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	
	rows, err := stmt.Query("%" + query + "%")
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

func (t *todo_mariadb) Close() error {
	return t.db.Close()
}
