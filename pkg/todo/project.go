package todo

import (
	"time"
)

// Project represents a project that can contain multiple todos
type Project struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"` // pointer to handle NULL in database
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ProjectService defines the interface for project management operations
type ProjectService interface {
	// CreateProject creates a new project with the given name and description
	CreateProject(name string, description *string) (Project, error)
	
	// GetAllProjects returns all projects
	GetAllProjects() []Project
	
	// GetProject returns a specific project by ID
	GetProject(id int64) (Project, error)
	
	// UpdateProject updates an existing project
	UpdateProject(id int64, name string, description *string) (Project, error)
	
	// DeleteProject deletes a project by ID
	DeleteProject(id int64) (Project, error)
	
	// GetProjectTodos returns all todos associated with a specific project
	GetProjectTodos(id int64) []TodoItem
}