package todo

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupSQLiteTestDB(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			completed BOOLEAN NOT NULL DEFAULT FALSE,
			due_date DATETIME,
			created_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	require.NoError(t, err)

	return db, func() {
		db.Close()
	}
}

func TestSQLiteTitleSearch(t *testing.T) {
	db, cleanup := setupSQLiteTestDB(t)
	defer cleanup()

	svc := NewTodoSQLite(db)

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

func TestSQLiteTitleSearch_WithSpecialCharacters(t *testing.T) {
	db, cleanup := setupSQLiteTestDB(t)
	defer cleanup()

	svc := NewTodoSQLite(db)

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
