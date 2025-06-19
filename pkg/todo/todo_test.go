package todo

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
)

func setupSQLiteTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			due_date DATETIME,
			created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	return db
}

func TestAddTodo(t *testing.T) {
	db := setupSQLiteTestDB(t)
	defer db.Close()

	svc := NewTodoSQLite(db)

func setupMariaDBTestDB(t *testing.T) *sql.DB {
	dsn := "root:password@tcp(localhost:3306)/testdb"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to connect to MariaDB: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			due_date DATETIME NULL,
			created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	// Clear any existing test data
	_, err = db.Exec("TRUNCATE TABLE todos")
	if err != nil {
		t.Fatalf("failed to truncate test table: %v", err)
	}

	return db
}

func TestMariaDBAddTodo(t *testing.T) {
	db := setupMariaDBTestDB(t)
	defer db.Close()

	svc := NewTodoMariaDB(db)

	// Same test cases as SQLite version

// Add similar test functions for all MariaDB operations following the same pattern as SQLite tests

	tests := []struct {
		name     string
		title    string
		dueDate  *time.Time
		wantErr  bool
	}{
		{
			name:    "valid todo",
			title:   "Test todo",
			dueDate: nil,
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			dueDate: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.AddTodo(tt.title, tt.dueDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompleteTodo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := NewTodoSQLite(db)

	// Add test todo
	item, err := svc.AddTodo("Test todo", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}

	// Test completing
	completed, err := svc.CompleteTodo(item.ID)
	if err != nil {
		t.Errorf("CompleteTodo() error = %v", err)
	}
	if !completed.Completed {
		t.Error("CompleteTodo() did not mark todo as completed")
	}

	// Test completing non-existent todo
	_, err = svc.CompleteTodo("nonexistent")
	if err == nil {
		t.Error("CompleteTodo() should fail for non-existent todo")
	}
}

func TestGetAllTodos(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := NewTodoSQLite(db)

	// Add test todos
	_, err := svc.AddTodo("Todo 1", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}
	_, err = svc.AddTodo("Todo 2", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}

	todos := svc.GetAllTodos()
	if len(todos) != 2 {
		t.Errorf("GetAllTodos() got %d todos, want 2", len(todos))
	}
}

func TestGetActiveCompletedTodos(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := NewTodoSQLite(db)

	// Add test todos
	item1, err := svc.AddTodo("Active 1", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}
	item2, err := svc.AddTodo("Active 2", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}

	// Complete one todo
	_, err = svc.CompleteTodo(item1.ID)
	if err != nil {
		t.Fatalf("failed to complete todo: %v", err)
	}

	// Test active todos
	active := svc.GetActiveTodos()
	if len(active) != 1 || active[0].ID != item2.ID {
		t.Error("GetActiveTodos() did not return expected active todos")
	}

	// Test completed todos
	completed := svc.GetCompletedTodos()
	if len(completed) != 1 || completed[0].ID != item1.ID {
		t.Error("GetCompletedTodos() did not return expected completed todos")
	}
}

func TestSetDueDate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := NewTodoSQLite(db)

	// Add test todo
	item, err := svc.AddTodo("Test todo", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}

	// Set due date
	dueDate := time.Now().Add(24 * time.Hour)
	updated, err := svc.SetDueDate(item.ID, dueDate)
	if err != nil {
		t.Errorf("SetDueDate() error = %v", err)
	}
	if updated.DueDate == nil || !updated.DueDate.Equal(dueDate) {
		t.Error("SetDueDate() did not update due date correctly")
	}
}

func TestTitleSearch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := NewTodoSQLite(db)

	// Add test todos
	_, err := svc.AddTodo("Buy groceries", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}
	_, err = svc.AddTodo("Clean house", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}

	// Test search
	results := svc.TitleSearchTodo("groceries")
	if len(results) != 1 || results[0].Title != "Buy groceries" {
		t.Error("TitleSearchTodo() did not return expected results")
	}
}

func TestDeleteTodo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	svc := NewTodoSQLite(db)

	// Add test todo
	item, err := svc.AddTodo("Test todo", nil)
	if err != nil {
		t.Fatalf("failed to add test todo: %v", err)
	}

	// Delete todo
	deleted, err := svc.DeleteTodo(item.ID)
	if err != nil {
		t.Errorf("DeleteTodo() error = %v", err)
	}
	if deleted.ID != item.ID {
		t.Error("DeleteTodo() did not return correct deleted item")
	}

	// Verify deleted
	_, err = svc.GetTodo(item.ID)
	if err == nil {
		t.Error("GetTodo() should fail for deleted todo")
	}
}
