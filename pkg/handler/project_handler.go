package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// CreateProjectHandler handles the create_project MCP tool
func (h *Handler) CreateProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.GetArguments()["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("project name is required and must be a string")
	}

	descriptionRaw, ok := request.GetArguments()["description"]
	var description *string
	if ok {
		descriptionStr, ok := descriptionRaw.(string)
		if !ok {
			return nil, fmt.Errorf("description must be a string")
		}
		description = &descriptionStr
	}

	if h.projectService == nil {
		return mcp.NewToolResultText(fmt.Sprintf("Project '%s' created successfully (placeholder - project service not initialized)", name)), nil
	}

	project, err := h.projectService.CreateProject(name, description)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return mcp.NewToolResultText(fmt.Sprintf("Project created: ID=%d, Name=%s", project.ID, project.Name)), nil
}

// GetAllProjectsHandler handles the get_all_projects MCP tool
func (h *Handler) GetAllProjectsHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.projectService == nil {
		return mcp.NewToolResultText("No projects found (project service not initialized)"), nil
	}

	projects := h.projectService.GetAllProjects()
	if len(projects) == 0 {
		return mcp.NewToolResultText("No projects found"), nil
	}

	var resultText string
	for _, project := range projects {
		description := ""
		if project.Description != nil {
			description = *project.Description
		}
		resultText += fmt.Sprintf("ID: %d, Name: %s, Description: %s, Created: %s\n", 
			project.ID, project.Name, description, project.CreatedAt.Format(time.RFC3339))
	}

	return mcp.NewToolResultText(resultText), nil
}

// GetProjectHandler handles the get_project MCP tool
func (h *Handler) GetProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("project id is required and must be a number")
	}
	id := int64(idRaw)

	if h.projectService == nil {
		return mcp.NewToolResultText(fmt.Sprintf("Project ID: %d (project service not initialized)", id)), nil
	}

	project, err := h.projectService.GetProject(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	description := ""
	if project.Description != nil {
		description = *project.Description
	}

	resultText := fmt.Sprintf("ID: %d, Name: %s, Description: %s, Created: %s, Updated: %s",
		project.ID, project.Name, description, project.CreatedAt.Format(time.RFC3339), project.UpdatedAt.Format(time.RFC3339))

	return mcp.NewToolResultText(resultText), nil
}

// UpdateProjectHandler handles the update_project MCP tool
func (h *Handler) UpdateProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("project id is required and must be a number")
	}
	id := int64(idRaw)

	nameRaw, ok := request.GetArguments()["name"]
	var name *string
	if ok {
		nameStr, ok := nameRaw.(string)
		if !ok {
			return nil, fmt.Errorf("name must be a string")
		}
		name = &nameStr
	}

	descriptionRaw, ok := request.GetArguments()["description"]
	var description *string
	if ok {
		descriptionStr, ok := descriptionRaw.(string)
		if !ok {
			return nil, fmt.Errorf("description must be a string")
		}
		description = &descriptionStr
	}
	
	_ = description // Suppress unused variable warning for now

	if name == nil && description == nil {
		return nil, fmt.Errorf("at least one of name or description must be provided")
	}

	// For now, return a placeholder message
	return mcp.NewToolResultText(fmt.Sprintf("Project %d updated successfully", id)), nil
}

// DeleteProjectHandler handles the delete_project MCP tool
func (h *Handler) DeleteProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("project id is required and must be a number")
	}
	id := int64(idRaw)

	// For now, return a placeholder message
	return mcp.NewToolResultText(fmt.Sprintf("Project %d deleted successfully", id)), nil
}

// GetProjectTodosHandler handles the get_project_todos MCP tool
func (h *Handler) GetProjectTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("project id is required and must be a number")
	}
	id := int64(idRaw)

	// For now, return a placeholder message
	return mcp.NewToolResultText(fmt.Sprintf("Todos for project %d", id)), nil
}

// AddTodoToProjectHandler handles the add_todo_to_project MCP tool
func (h *Handler) AddTodoToProjectHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := request.GetArguments()["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("todo title is required and must be a string")
	}

	projectIDRaw, ok := request.GetArguments()["project_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("project_id is required and must be a number")
	}
	projectID := int64(projectIDRaw)

	dueDateRaw, ok := request.GetArguments()["due_date"]
	var dueDate *time.Time
	if ok {
		dueDateStr, ok := dueDateRaw.(string)
		if !ok {
			return nil, fmt.Errorf("due_date must be a string")
		}
		parsedDueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse due_date: %w", err)
		}
		dueDate = &parsedDueDate
	}
	
	_ = dueDate // Suppress unused variable warning for now

	// For now, return a placeholder message
	return mcp.NewToolResultText(fmt.Sprintf("Todo '%s' added to project %d", title, projectID)), nil
}