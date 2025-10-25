package todo

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// NewCategorySQLite creates a new SQLite implementation of CategoryRepository
func NewCategorySQLite(db *sql.DB) CategoryRepository {
	return &category_sqlite{db: db}
}

type category_sqlite struct {
	db *sql.DB
}

// Create creates a new category with the given name, description, and color
func (c *category_sqlite) Create(category Category) (Category, error) {
	if category.Name == "" {
		return Category{}, fmt.Errorf("category name cannot be empty")
	}

	createdAt := time.Now()
	updatedAt := createdAt

	stmt, err := c.db.Prepare("INSERT INTO categories (name, description, color, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return Category{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(category.Name, category.Description, category.Color, createdAt, updatedAt)
	if err != nil {
		return Category{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Category{}, err
	}

	newCategory := Category{
		ID:          id,
		Name:        category.Name,
		Description: category.Description,
		Color:       category.Color,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	return newCategory, nil
}

// FindAll returns all categories
func (c *category_sqlite) FindAll() ([]Category, error) {
	stmt, err := c.db.Prepare("SELECT id, name, description, color, created_at, updated_at FROM categories ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.Name, &category.Description, &category.Color, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// FindByID returns a specific category by ID
func (c *category_sqlite) FindByID(id int64) (Category, error) {
	var category Category
	stmt, err := c.db.Prepare("SELECT id, name, description, color, created_at, updated_at FROM categories WHERE id = ?")
	if err != nil {
		return Category{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&category.ID, &category.Name, &category.Description, &category.Color, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Category{}, fmt.Errorf("category not found")
		}
		return Category{}, err
	}
	return category, nil
}

// FindByName returns a category by name
func (c *category_sqlite) FindByName(name string) (Category, error) {
	var category Category
	stmt, err := c.db.Prepare("SELECT id, name, description, color, created_at, updated_at FROM categories WHERE name = ?")
	if err != nil {
		return Category{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(name).Scan(&category.ID, &category.Name, &category.Description, &category.Color, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Category{}, fmt.Errorf("category not found")
		}
		return Category{}, err
	}
	return category, nil
}

// Update updates an existing category
func (c *category_sqlite) Update(category Category) (Category, error) {
	if category.Name == "" {
		return Category{}, fmt.Errorf("category name cannot be empty")
	}

	updatedAt := time.Now()

	stmt, err := c.db.Prepare("UPDATE categories SET name = ?, description = ?, color = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return Category{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(category.Name, category.Description, category.Color, updatedAt, category.ID)
	if err != nil {
		return Category{}, err
	}

	// Return the updated category
	return c.FindByID(category.ID)
}

// Delete deletes a category by ID
func (c *category_sqlite) Delete(id int64) error {
	stmt, err := c.db.Prepare("DELETE FROM categories WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

// FindTodosByCategory returns all todos associated with a specific category
func (c *category_sqlite) FindTodosByCategory(categoryID int64) ([]TodoItem, error) {
	stmt, err := c.db.Prepare("SELECT id, title, completed_at, due_date, created_date, reference_id, project_id, category_id FROM todos WHERE category_id = ? ORDER BY created_date DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []TodoItem
	for rows.Next() {
		var todo TodoItem
		err = rows.Scan(&todo.ID, &todo.Title, &todo.CompletedAt, &todo.DueDate, &todo.CreatedDate, &todo.ReferenceID, &todo.ProjectID, &todo.CategoryID)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// FindUncategorizedTodos returns all todos that are not assigned to any category
func (c *category_sqlite) FindUncategorizedTodos() ([]TodoItem, error) {
	stmt, err := c.db.Prepare("SELECT id, title, completed_at, due_date, created_date, reference_id, project_id, category_id FROM todos WHERE category_id IS NULL ORDER BY created_date DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []TodoItem
	for rows.Next() {
		var todo TodoItem
		err = rows.Scan(&todo.ID, &todo.Title, &todo.CompletedAt, &todo.DueDate, &todo.CreatedDate, &todo.ReferenceID, &todo.ProjectID, &todo.CategoryID)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}