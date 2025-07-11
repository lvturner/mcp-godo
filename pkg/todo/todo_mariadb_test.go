package todo

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mariadbTestDB *sql.DB

func TestMain(m *testing.M) {
	// Setup test database connection
	db, err := sql.Open("mysql", "root:password@tcp(localhost:3306)/testdb?parseTime=true")
	if err != nil {
		fmt.Printf("Failed to connect to test database: %v\n", err)
		os.Exit(1)
	}
	mariadbTestDB = db

	// Create recurrence_patterns table for tests
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS recurrence_patterns (
		id INT AUTO_INCREMENT PRIMARY KEY,
		todo_id VARCHAR(255) NOT NULL,
		frequency VARCHAR(50) NOT NULL,
		` + "`interval`" + ` INT NOT NULL,
		until DATETIME,
		count INT
	)`)
	if err != nil {
		fmt.Printf("Failed to create recurrence_patterns table: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	db.Exec("DELETE FROM recurrence_patterns")
	db.Exec("DELETE FROM todos")
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
	if todo.CompletedAt != nil {
		t.Error("New todo should not be completed")
	}
	if todo.CreatedDate.IsZero() || time.Since(todo.CreatedDate) > time.Second {
		t.Errorf("CreatedDate should be recent, got %v", todo.CreatedDate)
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
	if completed.CompletedAt == nil {
		t.Error("CompletedAt should be set after completing")
	} else {
		// Check that the completion time is recent (within the last minute)
		if time.Since(*completed.CompletedAt) > time.Minute {
			t.Errorf("CompletedAt is too old: %v", completed.CompletedAt)
		}
	}

	// Uncomplete it
	uncompleted, err := svc.UnCompleteTodo(todo.ID)
	if err != nil {
		t.Fatalf("UnCompleteTodo failed: %v", err)
	}
	if uncompleted.CompletedAt != nil {
		t.Error("CompletedAt should be nil after uncompleting")
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
	// TODO need extra check here - the due date doesn't come back EXACTLY as the date we entered
	// The returned due date just considers the date, the time is set to 00:00:00 
	if updated.DueDate == nil { 
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
		if todo.CompletedAt != nil {
			t.Error("Active todos should have nil CompletedAt")
		}
	}

	// Test GetCompletedTodos
	completedTodos := svc.GetCompletedTodos()
	if len(completedTodos) == 0 {
		t.Error("Expected completed todos")
	}
	for _, todo := range completedTodos {
		if todo.CompletedAt == nil {
			t.Error("Completed todos should have non-nil CompletedAt")
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

	// Setup test data
	testTodos := []struct {
		title     string
		completed bool
	}{
		{"Search active one", false},
		{"Search active two", false},
		{"Search COMPLETED one", true},
		{"Search completed two", true},
		{"Special !@#$%^&*() chars", false},
		{"Very long title " + strings.Repeat("a", 200), false},
		{"CaseSensitiveTest", false},
		{"", false}, // Empty title
	}

	// Add test todos
	var todoIDs []string
	for _, td := range testTodos {
		todo, err := svc.AddTodo(td.title, nil)
		if err != nil {
			if td.title == "" {
				continue // Expected to fail for empty title
			}
			t.Fatalf("AddTodo failed: %v", err)
		}
		todoIDs = append(todoIDs, todo.ID)
		if td.completed {
			_, err = svc.CompleteTodo(todo.ID)
			if err != nil {
				t.Fatalf("CompleteTodo failed: %v", err)
			}
		}
	}

	tests := []struct {
		name        string
		query       string
		activeOnly  bool
		expectedMin int
		validate    func([]TodoItem) error
	}{
		{
			name:        "General search",
			query:       "Search",
			activeOnly:  false,
			expectedMin: 4,
			validate: func(todos []TodoItem) error {
				if len(todos) < 4 {
					return fmt.Errorf("expected at least 4 todos, got %d", len(todos))
				}
				return nil
			},
		},
		{
			name:        "Active-only search",
			query:       "Search",
			activeOnly:  true,
			expectedMin: 2,
			validate: func(todos []TodoItem) error {
				for _, todo := range todos {
					if todo.CompletedAt != nil {
						return fmt.Errorf("found completed todo in active-only search")
					}
				}
				return nil
			},
		},
		{
			name:        "Case insensitive search",
			query:       "completed",
			activeOnly:  false,
			expectedMin: 2,
			validate: func(todos []TodoItem) error {
				if len(todos) < 2 {
					return fmt.Errorf("expected case insensitive match")
				}
				return nil
			},
		},
		{
			name:        "Special characters",
			query:       "!@#",
			activeOnly:  false,
			expectedMin: 1,
			validate: nil,
		},
		{
			name:        "Long title search",
			query:       strings.Repeat("a", 50),
			activeOnly:  false,
			expectedMin: 1,
			validate: nil,
		},
		{
			name:        "Empty query",
			query:       "",
			activeOnly:  false,
			expectedMin: len(testTodos) - 1, // minus the empty title one
			validate: nil,
		},
		{
			name:        "No matches",
			query:       "nonexistent",
			activeOnly:  false,
			expectedMin: 0,
			validate: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := svc.TitleSearchTodo(tt.query, tt.activeOnly)
			
			if len(results) < tt.expectedMin {
				t.Errorf("Expected at least %d results, got %d", tt.expectedMin, len(results))
			}
			
			if tt.validate != nil {
				if err := tt.validate(results); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
func TestMariaDB_RecurrencePattern(t *testing.T) {
	svc := NewTodoMariaDB(mariadbTestDB)

	// Add a test todo
	todo, err := svc.AddTodo("Recurring todo", nil)
	if err != nil {
		t.Fatalf("AddTodo failed: %v", err)
	}

	// Create recurrence pattern
	until := time.Now().AddDate(0, 1, 0) // 1 month from now
	count := 5
	pattern := RecurrencePattern{
		TodoID:    todo.ID,
		Frequency: "weekly",
		Interval:  1,
		Until:     &until,
		Count:     &count,
	}

	// Add recurrence pattern
	patternID, err := svc.AddRecurrencePattern(pattern)
	if err != nil {
		t.Fatalf("AddRecurrencePattern failed: %v", err)
	}
	if patternID <= 0 {
		t.Error("Invalid pattern ID returned")
	}

	// Retrieve recurrence pattern
	retrieved, err := svc.GetRecurrencePatternByID(patternID)
	if err != nil {
		t.Fatalf("GetRecurrencePatternByID failed: %v", err)
	}

	// Verify fields
	if retrieved.TodoID != todo.ID {
		t.Errorf("Expected TodoID %s, got %s", todo.ID, retrieved.TodoID)
	}
	if retrieved.Frequency != "weekly" {
		t.Errorf("Expected Frequency 'weekly', got '%s'", retrieved.Frequency)
	}
	if retrieved.Interval != 1 {
		t.Errorf("Expected Interval 1, got %d", retrieved.Interval)
	}
	if retrieved.Until == nil || !retrieved.Until.Truncate(24 * time.Hour).Equal(until.Truncate(24 * time.Hour)) {
		t.Error("Until date mismatch")
	}
	if retrieved.Count == nil || *retrieved.Count != count {
		t.Error("Count mismatch")
	}
}
