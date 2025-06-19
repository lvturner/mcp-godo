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
