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
	var dueDateStr sql.NullString
	err = t.db.QueryRow("SELECT title, completed, due_date FROM todos WHERE id = ?", id).Scan(
		&item.Title, &item.Completed, &dueDateStr)
	if err != nil {
		return TodoItem{}, err
	}
	
	// Parse due date if present - try both formats
	if dueDateStr.Valid {
		parsedDueDate, err := time.Parse(time.RFC3339, dueDateStr.String)
		if err != nil {
			// Try ISO 8601 format if MySQL format fails
			parsedDueDate, err = time.Parse(time.RFC3339, dueDateStr.String)
			if err != nil {
				return TodoItem{}, fmt.Errorf("error parsing due date: %w", err)
			}
		}
		item.DueDate = &parsedDueDate
	}
	
	return item, nil
}

func (t *todo_mariadb) CompleteTodo(id string) (TodoItem, error) {
	_, err := t.db.Exec("UPDATE todos SET completed = true WHERE id = ?", id)
	if err != nil {
		return TodoItem{}, err
	}
	item := TodoItem{ID: id}
	var dueDateStr sql.NullString
	err = t.db.QueryRow("SELECT title, completed, due_date FROM todos WHERE id = ?", id).Scan(
		&item.Title, &item.Completed, &dueDateStr)
	if dueDateStr.Valid {
		dueDate, err := time.Parse("2006-01-02 15:04:05", dueDateStr.String)
		if err != nil {
			return TodoItem{}, fmt.Errorf("error parsing due date: %w", err)
		}
		item.DueDate = &dueDate
	}
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
		var dueDateStr, createdDateStr sql.NullString
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &dueDateStr, &createdDateStr)
		if err != nil {
			log.Fatal(err)
		}
		
		// Parse created date - try both MySQL and ISO 8601 formats
		if createdDateStr.Valid {
			item.CreatedDate, err = time.Parse(time.RFC3339, createdDateStr.String)
			if err != nil {
				// Try ISO 8601 format if MySQL format fails
				item.CreatedDate, err = time.Parse(time.RFC3339, createdDateStr.String)
				if err != nil {
					log.Printf("error parsing created date: %v", err)
				}
			}
		}
		
		// Parse due date from MySQL format if present
		if dueDateStr.Valid {
			dueDate, err := time.Parse("2006-01-02 15:04:05", dueDateStr.String)
			if err != nil {
				return TodoItem{}, fmt.Errorf("error parsing due date: %w", err)
			}
			item.DueDate = &dueDate
		}
		items = append(items, item)
	}
	return items
}

func (t *todo_mariadb) GetTodo(id string) (TodoItem, error) {
	var item TodoItem
	var dueDateStr, createdDateStr sql.NullString
	err := t.db.QueryRow("SELECT id, title, completed, due_date, created_date FROM todos WHERE id = ?", id).Scan(
		&item.ID, &item.Title, &item.Completed, &dueDateStr, &createdDateStr)
	if err != nil {
		return TodoItem{}, err
	}
	
	// Parse created date - try both formats
	if createdDateStr.Valid {
		item.CreatedDate, err = time.Parse(time.RFC3339, createdDateStr.String)
		if err != nil {
			// Try ISO 8601 format if MySQL format fails
			item.CreatedDate, err = time.Parse(time.RFC3339, createdDateStr.String)
			if err != nil {
				return TodoItem{}, fmt.Errorf("error parsing created date: %w", err)
			}
		}
	}
	
	// Parse due date if present - try both formats
	if dueDateStr.Valid {
		dueDate, err := time.Parse(time.RFC3339, dueDateStr.String)
		if err != nil {
			// Try ISO 8601 format if MySQL format fails
			dueDate, err = time.Parse(time.RFC3339, dueDateStr.String)
			if err != nil {
				return TodoItem{}, fmt.Errorf("error parsing due date: %w", err)
			}
		}
		item.DueDate = &dueDate
	}
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
		var dueDateStr, createdDateStr sql.NullString
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &dueDateStr, &createdDateStr)
		if err != nil {
			log.Fatal(err)
		}
		
		// Parse created date
		if createdDateStr.Valid {
			item.CreatedDate, err = time.Parse(time.RFC3339, createdDateStr.String)
			if err != nil {
				log.Printf("error parsing created date: %v", err)
			}
		}
		
		// Parse due date if present
		if dueDateStr.Valid {
			dueDate, err := time.Parse("2006-01-02 15:04:05", dueDateStr.String)
			if err != nil {
				log.Printf("error parsing due date: %v", err)
			} else {
				item.DueDate = &dueDate
			}
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
		var dueDateStr, createdDateStr sql.NullString
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &dueDateStr, &createdDateStr)
		if err != nil {
			log.Fatal(err)
		}
		
		// Parse created date
		if createdDateStr.Valid {
			item.CreatedDate, err = time.Parse(time.RFC3339, createdDateStr.String)
			if err != nil {
				log.Printf("error parsing created date: %v", err)
			}
		}
		
		// Parse due date if present
		if dueDateStr.Valid {
			dueDate, err := time.Parse(time.RFC3339, dueDateStr.String)
			if err != nil {
				log.Printf("error parsing due date: %v", err)
			} else {
				item.DueDate = &dueDate
			}
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

func (t *todo_mariadb) TitleSearchTodo(query string) []TodoItem {
	if query == "" {
		return t.GetAllTodos()
	}

	// Use prepared statement to prevent SQL injection
	stmt, err := t.db.Prepare(`
		SELECT id, title, completed, due_date, created_date 
		FROM todos 
		WHERE title LIKE ?
	`)
	if err != nil {
		log.Printf("error preparing search statement: %v", err)
		return nil
	}
	defer stmt.Close()

	rows, err := stmt.Query("%" + query + "%")
	if err != nil {
		log.Printf("error searching todos: %v", err)
		return nil
	}
	defer rows.Close()
	
	var items []TodoItem
	for rows.Next() {
		var item TodoItem
		var dueDateStr, createdDateStr sql.NullString
		err = rows.Scan(&item.ID, &item.Title, &item.Completed, &dueDateStr, &createdDateStr)
		if err != nil {
			log.Printf("error scanning todo row: %v", err)
			continue
		}
			
		// Parse created date
		if createdDateStr.Valid {
			item.CreatedDate, err = time.Parse(time.RFC3339, createdDateStr.String)
			if err != nil {
				log.Printf("error parsing created date: %v", err)
				continue
			}
		}
			
		// Parse due date if present
		if dueDateStr.Valid {
			dueDate, err := time.Parse(time.RFC3339, dueDateStr.String)
			if err != nil {
				log.Printf("error parsing due date: %v", err)
			} else {
				item.DueDate = &dueDate
			}
		}
		items = append(items, item)
	}
	
	if err = rows.Err(); err != nil {
		log.Printf("error after scanning rows: %v", err)
	}
	
	return items
}

func (t *todo_mariadb) Close() error {
	return t.db.Close()
}
