package todo

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMariaDBTestDB(t *testing.T) (*sql.DB, func()) {
	dsn := os.Getenv("MARIADB_TEST_DSN")
	if dsn == "" {
		t.Skip("MARIADB_TEST_DSN not set, skipping MariaDB tests")
	}

	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)

	// Create test table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTO_INCREMENT,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			due_date DATETIME,
			created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	// Clear any existing data
	_, err = db.Exec("DELETE FROM todos")
	require.NoError(t, err)

	return db, func() {
		_, _ = db.Exec("DROP TABLE todos")
		db.Close()
	}
}

func TestMariaDBBasicOperations(t *testing.T) {
	db, cleanup := setupMariaDBTestDB(t)
	defer cleanup()

	svc := NewTodoMariaDB(db)

	// Test AddTodo
	item, err := svc.AddTodo("Test todo", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Test todo", item.Title)
	assert.False(t, item.Completed)
	assert.NotEmpty(t, item.ID)

	// Test GetTodo
	retrieved, err := svc.GetTodo(item.ID)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, retrieved.ID)
	assert.Equal(t, item.Title, retrieved.Title)

	// Test CompleteTodo
	completed, err := svc.CompleteTodo(item.ID)
	assert.NoError(t, err)
	assert.True(t, completed.Completed)

	// Test UnCompleteTodo
	uncompleted, err := svc.UnCompleteTodo(item.ID)
	assert.NoError(t, err)
	assert.False(t, uncompleted.Completed)

	// Test DeleteTodo
	deleted, err := svc.DeleteTodo(item.ID)
	assert.NoError(t, err)
	assert.Equal(t, item.ID, deleted.ID)

	// Verify deleted
	_, err = svc.GetTodo(item.ID)
	assert.Error(t, err)
}

func TestMariaDBWithDueDate(t *testing.T) {
	db, cleanup := setupMariaDBTestDB(t)
	defer cleanup()

	svc := NewTodoMariaDB(db)

	dueDate := time.Now().Add(24 * time.Hour)
	item, err := svc.AddTodo("Todo with due date", &dueDate)
	assert.NoError(t, err)
	assert.NotNil(t, item.DueDate)
	assert.Equal(t, dueDate.Truncate(time.Second), item.DueDate.Truncate(time.Second))

	// Test SetDueDate
	newDueDate := time.Now().Add(48 * time.Hour)
	updated, err := svc.SetDueDate(item.ID, newDueDate)
	assert.NoError(t, err)
	assert.Equal(t, newDueDate.Truncate(time.Second), updated.DueDate.Truncate(time.Second))
}

func TestMariaDBGetAllActiveCompleted(t *testing.T) {
	db, cleanup := setupMariaDBTestDB(t)
	defer cleanup()

	svc := NewTodoMariaDB(db)

	// Add test data
	titles := []string{"Active 1", "Active 2", "Completed 1", "Completed 2"}
	for i, title := range titles {
		item, err := svc.AddTodo(title, nil)
		assert.NoError(t, err)
		if i >= 2 {
			_, err = svc.CompleteTodo(item.ID)
			assert.NoError(t, err)
		}
	}

	// Test GetAllTodos
	all := svc.GetAllTodos()
	assert.Len(t, all, 4)

	// Test GetActiveTodos
	active := svc.GetActiveTodos()
	assert.Len(t, active, 2)
	for _, item := range active {
		assert.False(t, item.Completed)
	}

	// Test GetCompletedTodos
	completed := svc.GetCompletedTodos()
	assert.Len(t, completed, 2)
	for _, item := range completed {
		assert.True(t, item.Completed)
	}
}

func TestMariaDBTitleSearch(t *testing.T) {
	db, cleanup := setupMariaDBTestDB(t)
	defer cleanup()

	svc := NewTodoMariaDB(db)

	// Add test data
	testTitles := []string{
		"Buy groceries",
		"Finish project report",
		"Buy new laptop",
		"Call mom",
		"Buy birthday gift",
		"Project meeting",
	}
	for _, title := range testTitles {
		_, err := svc.AddTodo(title, nil)
		assert.NoError(t, err)
	}

	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "exact match",
			query:    "Buy groceries",
			expected: []string{"Buy groceries"},
		},
		{
			name:     "partial match",
			query:    "Buy",
			expected: []string{"Buy groceries", "Buy new laptop", "Buy birthday gift"},
		},
		{
			name:     "case insensitive",
			query:    "buy",
			expected: []string{"Buy groceries", "Buy new laptop", "Buy birthday gift"},
		},
		{
			name:     "partial word match",
			query:    "proj",
			expected: []string{"Finish project report", "Project meeting"},
		},
		{
			name:     "no match",
			query:    "nonexistent",
			expected: []string{},
		},
		{
			name:     "empty query returns all",
			query:    "",
			expected: testTitles,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := svc.TitleSearchTodo(tt.query)
			assert.Equal(t, len(tt.expected), len(results), "unexpected number of results")
			
			resultTitles := make([]string, len(results))
			for i, item := range results {
				resultTitles[i] = item.Title
			}
			
			for _, expected := range tt.expected {
				assert.Contains(t, resultTitles, expected, "expected title not found in results")
			}
		})
	}
}

func TestMariaDBTitleSearch_WithSpecialCharacters(t *testing.T) {
	db, cleanup := setupMariaDBTestDB(t)
	defer cleanup()

	svc := NewTodoMariaDB(db)

	specialTitles := []string{
		"Meeting @ 2pm",
		"Email client re: project",
		"Fix bug #1234",
		"Update docs (urgent)",
	}
	for _, title := range specialTitles {
		_, err := svc.AddTodo(title, nil)
		assert.NoError(t, err)
	}

	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name:     "search with @",
			query:    "@",
			expected: []string{"Meeting @ 2pm"},
		},
		{
			name:     "search with :",
			query:    "re:",
			expected: []string{"Email client re: project"},
		},
		{
			name:     "search with #",
			query:    "#1234",
			expected: []string{"Fix bug #1234"},
		},
		{
			name:     "search with parentheses",
			query:    "(urgent)",
			expected: []string{"Update docs (urgent)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := svc.TitleSearchTodo(tt.query)
			assert.Equal(t, len(tt.expected), len(results))
			if len(results) > 0 {
				assert.Equal(t, tt.expected[0], results[0].Title)
			}
		})
	}
}
