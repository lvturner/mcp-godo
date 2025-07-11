package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"mcp-godo/pkg/todo"

	"github.com/mark3labs/mcp-go/mcp"
)

type Handler struct {
	todoService todo.TodoService
}

func NewHandler(todoService todo.TodoService) *Handler {
	return &Handler{todoService: todoService}
}

func (h *Handler) AddRecurrencePatternHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todoID, ok := request.GetArguments()["todo_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid todo_id")
	}
	frequency, ok := request.GetArguments()["frequency"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid frequency")
	}
	interval, ok := request.GetArguments()["interval"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid interval")
	}
	untilRaw, ok := request.GetArguments()["until"]
	var until *time.Time
	if ok {
		untilStr, ok := untilRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid until")
		}
		untilTime, err := time.Parse(time.RFC3339, untilStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse until: %w", err)
		}
		until = &untilTime
	}
	countRaw, ok := request.GetArguments()["count"]
	var count *int
	if ok {
		countVal, ok := countRaw.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid count")
		}
		countInt := int(countVal)
		count = &countInt
	}

	pattern := todo.RecurrencePattern{
		TodoID:    todoID,
		Frequency: frequency,
		Interval:  int(interval),
		Until:     until,
		Count:     count,
	}

	patternID, err := h.todoService.AddRecurrencePattern(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to add recurrence pattern: %w", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Recurrence pattern added with ID: %d", patternID)), nil
}

func (h *Handler) GetRecurrencePatternHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	idRaw, ok := request.GetArguments()["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid id")
	}
	id := int64(idRaw)
	pattern, err := h.todoService.GetRecurrencePatternByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get recurrence pattern: %w", err)
	}
	resultText := fmt.Sprintf("ID: %d, TodoID: %s, Frequency: %s, Interval: %d", 
		pattern.ID, pattern.TodoID, pattern.Frequency, pattern.Interval)
	if pattern.Until != nil {
		resultText += fmt.Sprintf(", Until: %s", pattern.Until.Format(time.RFC3339))
	}
	if pattern.Count != nil {
		resultText += fmt.Sprintf(", Count: %d", *pattern.Count)
	}
	return mcp.NewToolResultText(resultText), nil
}

func (h *Handler) TitleSearchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.GetArguments()["query"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid query")
	}
	activeOnly, _ := request.GetArguments()["active_only"].(bool)
	todos := h.todoService.TitleSearchTodo(query, activeOnly)
	var results []string
	for _, todo := range todos {
		var dueDateStr string
		if todo.DueDate != nil {
			dueDateStr = todo.DueDate.Format(time.RFC3339)
		}
		referenceID := ""
		if todo.ReferenceID != nil {
			referenceID = fmt.Sprintf("ReferenceID: %d", *todo.ReferenceID)
		}
		results = append(results, fmt.Sprintf("ID: %s, Title: %s, CompletedAt: %s, Due Date: %s, %s", 
			todo.ID, todo.Title, todo.CompletedAt, dueDateStr, referenceID))
	}

	if len(results) > 0 {
		return mcp.NewToolResultText(strings.Join(results, "\n")), nil
	} else {
		return mcp.NewToolResultText("No todos found"), nil
	}
}

func (h *Handler) UpdateDueDateHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok{
		return nil, fmt.Errorf("invalid id")
	}
	dueDateStr, ok := request.GetArguments()["due_date"].(string)
	if !ok{
		return nil, fmt.Errorf("invalid due date")
	}
	dueDate, err := time.Parse(time.RFC3339, dueDateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse due date: %w", err)
	}
	todo, err := h.todoService.SetDueDate(id, dueDate)
	if err != nil{
		return nil, fmt.Errorf("failed to update due date: %w", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Todo updated: ID=%s, Title=%s, Due Date=%s", todo.ID, todo.Title, todo.DueDate.Format(time.RFC3339))), nil
}

func (h *Handler) UnCompleteTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok{
		return nil, fmt.Errorf("invalid id")
	}
	todo, err := h.todoService.UnCompleteTodo(id)
	if err != nil{
		return nil, fmt.Errorf("failed to uncomplete todo: %w", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Todo uncompleted: ID=%s, Title=%s", todo.ID, todo.Title)), nil
}

func (h *Handler) GetCompletedTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todos := h.todoService.GetCompletedTodos()
	if len(todos) == 0 {
		return mcp.NewToolResultText("No completed todos found"), nil
	}
	var resultText string
	for _, todo := range todos {
		status := "Completed"
		referenceID := ""
		if todo.ReferenceID != nil {
			referenceID = fmt.Sprintf(", ReferenceID: %d", *todo.ReferenceID)
		}
		resultText += fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s%s\n", 
			todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate, referenceID)
	}
	return mcp.NewToolResultText(resultText), nil
}

func (h *Handler) GetActiveTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todos := h.todoService.GetActiveTodos()
	if len(todos) == 0 {
		return mcp.NewToolResultText("No active todos found"), nil
	}
	var resultText string
	for _, todo := range todos {
		status := "Incomplete"
		referenceID := ""
		if todo.ReferenceID != nil {
			referenceID = fmt.Sprintf(", ReferenceID: %d", *todo.ReferenceID)
		}
		resultText += fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s%s\n", 
			todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate, referenceID)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Today's date is %s, and the list of todo items is: %s", time.Now().Format("2006-01-02"), resultText)), nil
}

func (h *Handler) DeleteTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok {
		return nil, errors.New("id must be a string")
	}
	todo, err := h.todoService.DeleteTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete todo: %w", err)	
	}

	resultText := fmt.Sprintf("Deleted Todo: ID=%s, Title=%s", todo.ID, todo.Title)
	return mcp.NewToolResultText(resultText), nil
}

func (h *Handler) GetTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok {
		return nil, errors.New("id must be a string")
	}
	todo, err := h.todoService.GetTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	status := "Incomplete"
	if todo.CompletedAt != nil {
		status = "Complete"
	}
	referenceID := ""
	if todo.ReferenceID != nil {
		referenceID = fmt.Sprintf(", ReferenceID: %d", *todo.ReferenceID)
	}

	resultText := fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s%s\n", 
		todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate, referenceID)
	
	return mcp.NewToolResultText(resultText), nil
}

func (h *Handler) ListTodosHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	todos := h.todoService.GetAllTodos()
	var todosText []string
	for _, todo := range todos {
		status := "Incomplete"
		if todo.CompletedAt != nil {
			status = "Complete"
		}
		referenceID := ""
		if todo.ReferenceID != nil {
			referenceID = fmt.Sprintf(", ReferenceID: %d", *todo.ReferenceID)
		}
		todosText = append(todosText, fmt.Sprintf("ID: %s, Title: %s, Status: %s, Due Date: %s, Created Date: %s%s\n", 
			todo.ID, todo.Title, status, todo.DueDate, todo.CreatedDate, referenceID))
	}
	return mcp.NewToolResultText(strings.Join(todosText, "\n")), nil
}

func (h *Handler) CompleteTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, ok := request.GetArguments()["id"].(string)
	if !ok {
		return nil, errors.New("id must be a string")
	}

	completedTodo, err := h.todoService.CompleteTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to complete todo: %v", err)
	}
	return mcp.NewToolResultText(fmt.Sprintf("Todo %s completed", completedTodo.Title)), nil
}

func (h *Handler) AddTodoHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, ok := request.GetArguments()["title"].(string)
	if !ok {
		return nil, errors.New("title must be a string")
	}
	dueDateRaw, ok := request.GetArguments()["due_date"]
	if ok {
		dueDateStr, ok := dueDateRaw.(string)
		if !ok {
			return nil, errors.New("due_date must be a string")
		}
		dueDate, err := time.Parse(time.RFC3339, dueDateStr)
		if err != nil {
			return nil, err
		}

		h.todoService.AddTodo(title, &dueDate)
	} else {
		h.todoService.AddTodo(title, nil)
	}
	
	return mcp.NewToolResultText(fmt.Sprintf("%s added to todo list", title)), nil
}

func (h *Handler) ListTodosResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	todos := h.todoService.GetAllTodos()
	
	jsonData, err := json.Marshal(todos)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal todos to JSON: %w", err)
	}
	
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI: request.Params.URI,
			MIMEType: "application/json",
			Text: string(jsonData),
		},
	}, nil
}

func (h *Handler) GetSingleTodoResourceHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	id := extractIDFromURI(request.Params.URI)
	todo, err := h.todoService.GetTodo(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo item: %w", err)
	}	
	
	jsonData, err := json.Marshal(todo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal todo item: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:	request.Params.URI,
			MIMEType: "application/json",
			Text: string(jsonData),
		},
	}, nil
}

func extractIDFromURI(uri string) string {
    parsed, err := url.Parse(uri)
    if err != nil || parsed.Scheme != "todos" || parsed.Opaque == "" {
        return ""
    }

    return strings.TrimPrefix(parsed.Path, "/") 
}
