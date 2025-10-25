package todo

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// NewProjectMariaDB creates a new MySQL implementation of ProjectService
func NewProjectMariaDB(db *sql.DB) ProjectService {
	return &project_mariadb{db: db}
}

type project_mariadb struct {
	db *sql.DB
}

// CreateProject creates a new project with the given name and description
func (p *project_mariadb) CreateProject(name string, description *string) (Project, error) {
	if name == "" {
		return Project{}, fmt.Errorf("project name cannot be empty")
	}

	createdAt := time.Now()
	updatedAt := createdAt

	stmt, err := p.db.Prepare("INSERT INTO projects (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return Project{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, description, createdAt, updatedAt)
	if err != nil {
		return Project{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return Project{}, err
	}

	newProject := Project{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	return newProject, nil
}

// GetAllProjects returns all projects
func (p *project_mariadb) GetAllProjects() []Project {
	stmt, err := p.db.Prepare("SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		err = rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}
		projects = append(projects, project)
	}
	return projects
}

// GetProject returns a specific project by ID
func (p *project_mariadb) GetProject(id int64) (Project, error) {
	var project Project
	stmt, err := p.db.Prepare("SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ?")
	if err != nil {
		return Project{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Project{}, fmt.Errorf("project not found")
		}
		return Project{}, err
	}
	return project, nil
}

// UpdateProject updates an existing project
func (p *project_mariadb) UpdateProject(id int64, name string, description *string) (Project, error) {
	if name == "" {
		return Project{}, fmt.Errorf("project name cannot be empty")
	}

	updatedAt := time.Now()

	stmt, err := p.db.Prepare("UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return Project{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description, updatedAt, id)
	if err != nil {
		return Project{}, err
	}

	// Return the updated project
	return p.GetProject(id)
}

// DeleteProject deletes a project by ID
func (p *project_mariadb) DeleteProject(id int64) (Project, error) {
	// First get the project to return it after deletion
	project, err := p.GetProject(id)
	if err != nil {
		return Project{}, err
	}

	stmt, err := p.db.Prepare("DELETE FROM projects WHERE id = ?")
	if err != nil {
		return Project{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return Project{}, err
	}

	return project, nil
}

// GetProjectTodos returns all todos associated with a specific project
func (p *project_mariadb) GetProjectTodos(id int64) []TodoItem {
	stmt, err := p.db.Prepare("SELECT id, title, completed_at, due_date, created_date, reference_id, project_id FROM todos WHERE project_id = ? ORDER BY created_date DESC")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var todos []TodoItem
	for rows.Next() {
		var todo TodoItem
		err = rows.Scan(&todo.ID, &todo.Title, &todo.CompletedAt, &todo.DueDate, &todo.CreatedDate, &todo.ReferenceID, &todo.ProjectID)
		if err != nil {
			log.Fatal(err)
		}
		todos = append(todos, todo)
	}
	return todos
}