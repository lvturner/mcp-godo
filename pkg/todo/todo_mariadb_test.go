package todo

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mariadbTestDB *sql.DB

func TestMain(m *testing.M) {
	// Setup test database connection
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb")
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}
	mariadbTestDB = db

	// Run tests
	code := m.Run()

	// Cleanup
	db.Close()
	os.Exit(code)
}

func TestMariaDB_AddTodo(t *testing.T) {
	svc := NewTodoMariaDB(mariadbTestDB)

	// Test adding todo without due date
	todo, err := svc.AddTodo("Test todo", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}
	if todo.Title != "Test todo" {
		t.Errorf("Expected title 'Test todo', got '%s'", todo.Title)
	}
	if todo.Completed {
		t.Error("New todo should not be completed")
	}

	// Test adding todo with due date
	dueDate := time.Now().Add(24 * time.Hour)
	todoWithDate, err := svc.AddTodo("Dated todo", &dueDate)
	if err != nil {
		t.Fatalf("AddTodo with due date failed: %v", err)
	}
	if todoWithDate.DueDate == nil || !todoWithDate.DueDate.Equal(dueDate) {
		t.Error("Due date not set correctly")
	}

	// Test empty title
	_, err = svc.AddTodo("", nil)
	if err == nil {
		t.Error("Expected error for empty title")
	}
}

func TestMariaDB_CompleteUncomplete(t *testing.T) {
	svc := NewTodoMariaDB(mariadbTestDB)

	// Add test todo
	todo, err := svc.AddTodo("Complete test", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}

	// Complete it
	completed, err := svc.CompleteTodo(todo.ID)
	if err != nil {
		t.Fatalf("CompleteTodo failed: %v", err)
	}
	if !completed.Completed {
		t.Error("Todo should be completed")
	}

	// Uncomplete it
	uncompleted, err := svc.UnCompleteTodo(todo.ID)
	if err != nil {
		t.Fatalf("UnCompleteTodo failed: %v", err)
	}
	if uncompleted.Completed {
		t.Error("Todo should be uncompleted")
	}
}

func TestMariaDB_SetDueDate(t *testing.T) {
	svc := NewTodoMariaDB(mariadbTestDB)

	// Add test todo
	todo, err := svc.AddTodo("Due date test", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}

	// Set due date
	newDate := time.Now().Add(48 * time.Hour)
	updated, err := svc.SetDueDate(todo.ID, newDate)
	if err != nil {
		t.Fatalf("SetDueDate failed: %v", err)
	}
	if updated.DueDate == nil || !updated.DueDate.Equal(newDate) {
		t.Error("Due date not updated correctly")
	}
}

func TestMariaDB_GetOperations(t *testing.T) {
	svc := NewTodoMariaDB(mariadbTestDB)

	// Add test todos
	_, err := svc.AddTodo("Active 1", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}
	completed, err := svc.AddTodo("Completed 1", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}
	_, err = svc.CompleteTodo(completed.ID)
	if err != nil {
		t.Fatalf("CompleteTodo failed: %v", err)
	}

	// Test GetAllTodos
	all := svc.GetAllTodos()
	if len(all) < 2 {
		t.Errorf("Expected at least 2 todos, got %d", len(all))
	}

	// Test GetActiveTodos
	active := svc.GetActiveTodos()
	if len(active) == 0 {
		t.Error("Expected active todos")
	}
	for _, todo := range active {
		if todo.Completed {
			t.Error("Active todos should not be completed")
		}
	}

	// Test GetCompletedTodos
	completedTodos := svc.GetCompletedTodos()
	if len(completedTodos) == 0 {
		t.Error("Expected completed todos")
	}
	for _, todo := range completedTodos {
		if !todo.Completed {
			t.Error("Completed todos should be marked as such")
		}
	}

	// Test GetTodo
	fetched, err := svc.GetTodo(completed.ID)
	if err != nil {
		t.Fatalf("GetTodo failed: %v", err)
	}
	if fetched.ID != completed.ID {
		t.Error("Fetched todo ID mismatch")
	}
}

func TestMariaDB_DeleteTodo(t *testing.T) {
	svc := NewTodoMariaDB(mariadbTestDB)

	// Add test todo
	todo, err := svc.AddTodo("Delete test", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}

	// Delete it
	deleted, err := svc.DeleteTodo(todo.ID)
	if err != nil {
		t.Fatalf("DeleteTodo failed: %v", err)
	}
	if deleted.ID != todo.ID {
		t.Error("Deleted todo ID mismatch")
	}

	// Verify it's gone
	_, err = svc.GetTodo(todo.ID)
	if err == nil {
		t.Error("Expected error when fetching deleted todo")
	}
}

func TestMariaDB_TitleSearch(t *testing.T) {
	svc := NewTodoMariaDB(mariadbTestDB)

	// Add test todos
	_, err := svc.AddTodo("Search test one", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}
	_, err = svc.AddTodo("Search test two", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}
	_, err = svc.AddTodo("Other item", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}

	// Search
	results := svc.TitleSearchTodo("Search test")
	if len(results) != 2 {
		t.Errorf("Expected 2 search results, got %d", len(results))
	}
}
