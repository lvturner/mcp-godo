package todo

import (
	"fmt"
	"time"
)

// Category represents a grouping of todo items by type or theme
type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"` // Optional description
	Color       *string   `json:"color"`       // Optional hex color code (e.g., "#FF5733")
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryService defines the interface for category operations
type CategoryService interface {
	// Category CRUD operations
	CreateCategory(name string, description *string, color *string) (Category, error)
	GetAllCategories() ([]Category, error)
	GetCategoryByID(id int64) (Category, error)
	UpdateCategory(id int64, name *string, description *string, color *string) (Category, error)
	DeleteCategory(id int64) error
	
	// Category-todo relationship operations
	GetTodosByCategory(categoryID int64) ([]TodoItem, error)
	GetUncategorizedTodos() ([]TodoItem, error)
}

// CategoryRepository defines the database operations for categories
type CategoryRepository interface {
	// Category CRUD operations
	Create(category Category) (Category, error)
	FindAll() ([]Category, error)
	FindByID(id int64) (Category, error)
	FindByName(name string) (Category, error)
	Update(category Category) (Category, error)
	Delete(id int64) error
	
	// Category-todo relationship operations
	FindTodosByCategory(categoryID int64) ([]TodoItem, error)
	FindUncategorizedTodos() ([]TodoItem, error)
}

// categoryService implements the CategoryService interface
type categoryService struct {
	repo CategoryRepository
}

// NewCategoryService creates a new category service instance
func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

// CreateCategory creates a new category with validation
func (s *categoryService) CreateCategory(name string, description *string, color *string) (Category, error) {
	// Validate name
	if name == "" {
		return Category{}, fmt.Errorf("category name cannot be empty")
	}
	
	// Validate name length
	if len(name) > 255 {
		return Category{}, fmt.Errorf("category name cannot exceed 255 characters")
	}
	
	// Validate color format if provided
	if color != nil && *color != "" {
		if !isValidHexColor(*color) {
			return Category{}, fmt.Errorf("invalid hex color format")
		}
	}
	
	// Check for duplicate name
	existing, err := s.repo.FindByName(name)
	if err == nil && existing.ID != 0 {
		return Category{}, fmt.Errorf("category with name '%s' already exists", name)
	}
	
	category := Category{
		Name:        name,
		Description: description,
		Color:       color,
	}
	
	return s.repo.Create(category)
}

// GetAllCategories retrieves all categories
func (s *categoryService) GetAllCategories() ([]Category, error) {
	return s.repo.FindAll()
}

// GetCategoryByID retrieves a category by ID
func (s *categoryService) GetCategoryByID(id int64) (Category, error) {
	return s.repo.FindByID(id)
}

// UpdateCategory updates an existing category with validation
func (s *categoryService) UpdateCategory(id int64, name *string, description *string, color *string) (Category, error) {
	// Get existing category
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return Category{}, err
	}
	
	// Update fields if provided
	if name != nil {
		if *name == "" {
			return Category{}, fmt.Errorf("category name cannot be empty")
		}
		if len(*name) > 255 {
			return Category{}, fmt.Errorf("category name cannot exceed 255 characters")
		}
		existing.Name = *name
	}
	
	if description != nil {
		existing.Description = description
	}
	
	if color != nil {
		if *color != "" && !isValidHexColor(*color) {
			return Category{}, fmt.Errorf("invalid hex color format")
		}
		existing.Color = color
	}
	
	return s.repo.Update(existing)
}

// DeleteCategory deletes a category by ID
func (s *categoryService) DeleteCategory(id int64) error {
	return s.repo.Delete(id)
}

// GetTodosByCategory retrieves all todos assigned to a specific category
func (s *categoryService) GetTodosByCategory(categoryID int64) ([]TodoItem, error) {
	return s.repo.FindTodosByCategory(categoryID)
}

// GetUncategorizedTodos retrieves all todos without a category
func (s *categoryService) GetUncategorizedTodos() ([]TodoItem, error) {
	return s.repo.FindUncategorizedTodos()
}

// isValidHexColor validates hex color format
func isValidHexColor(color string) bool {
	if len(color) != 7 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		if !((color[i] >= '0' && color[i] <= '9') || (color[i] >= 'a' && color[i] <= 'f') || (color[i] >= 'A' && color[i] <= 'F')) {
			return false
		}
	}
	return true
}