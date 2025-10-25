package todo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProjectService is a mock implementation of ProjectService for testing
type MockProjectService struct {
	mock.Mock
}

func (m *MockProjectService) CreateProject(name string, description *string) (Project, error) {
	args := m.Called(name, description)
	return args.Get(0).(Project), args.Error(1)
}

func (m *MockProjectService) GetAllProjects() []Project {
	args := m.Called()
	return args.Get(0).([]Project)
}

func (m *MockProjectService) GetProject(id int64) (Project, error) {
	args := m.Called(id)
	return args.Get(0).(Project), args.Error(1)
}

func (m *MockProjectService) UpdateProject(id int64, name string, description *string) (Project, error) {
	args := m.Called(id, name, description)
	return args.Get(0).(Project), args.Error(1)
}

func (m *MockProjectService) DeleteProject(id int64) (Project, error) {
	args := m.Called(id)
	return args.Get(0).(Project), args.Error(1)
}

func (m *MockProjectService) GetProjectTodos(id int64) []TodoItem {
	args := m.Called(id)
	return args.Get(0).([]TodoItem)
}

func TestProjectModel(t *testing.T) {
	// Test Project struct creation
	now := time.Now()
	description := "Test project description"
	
	project := Project{
		ID:          1,
		Name:        "Test Project",
		Description: &description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.Equal(t, int64(1), project.ID)
	assert.Equal(t, "Test Project", project.Name)
	assert.Equal(t, &description, project.Description)
	assert.Equal(t, now, project.CreatedAt)
	assert.Equal(t, now, project.UpdatedAt)
}

func TestProjectServiceInterface(t *testing.T) {
	// Test that MockProjectService implements ProjectService interface
	var _ ProjectService = (*MockProjectService)(nil)
}

func TestProjectServiceMock(t *testing.T) {
	mockService := new(MockProjectService)
	
	// Test CreateProject
	description := "Test description"
	expectedProject := Project{
		ID:          1,
		Name:        "Test Project",
		Description: &description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	mockService.On("CreateProject", "Test Project", &description).Return(expectedProject, nil)
	
	project, err := mockService.CreateProject("Test Project", &description)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedProject.ID, project.ID)
	assert.Equal(t, expectedProject.Name, project.Name)
	assert.Equal(t, expectedProject.Description, project.Description)
	
	mockService.AssertExpectations(t)
}

func TestProjectServiceGetAllProjects(t *testing.T) {
	mockService := new(MockProjectService)
	
	expectedProjects := []Project{
		{
			ID:        1,
			Name:      "Project 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Project 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	
	mockService.On("GetAllProjects").Return(expectedProjects)
	
	projects := mockService.GetAllProjects()
	
	assert.Len(t, projects, 2)
	assert.Equal(t, "Project 1", projects[0].Name)
	assert.Equal(t, "Project 2", projects[1].Name)
	
	mockService.AssertExpectations(t)
}

func TestProjectServiceGetProject(t *testing.T) {
	mockService := new(MockProjectService)
	
	description := "Test description"
	expectedProject := Project{
		ID:          1,
		Name:        "Test Project",
		Description: &description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	mockService.On("GetProject", int64(1)).Return(expectedProject, nil)
	
	project, err := mockService.GetProject(1)
	
	assert.NoError(t, err)
	assert.Equal(t, expectedProject.ID, project.ID)
	assert.Equal(t, expectedProject.Name, project.Name)
	
	mockService.AssertExpectations(t)
}

func TestProjectServiceGetProjectTodos(t *testing.T) {
	mockService := new(MockProjectService)
	
	expectedTodos := []TodoItem{
		{
			ID:        "1",
			Title:     "Todo 1",
			CreatedDate: time.Now(),
		},
		{
			ID:        "2",
			Title:     "Todo 2",
			CreatedDate: time.Now(),
		},
	}
	
	mockService.On("GetProjectTodos", int64(1)).Return(expectedTodos)
	
	todos := mockService.GetProjectTodos(1)
	
	assert.Len(t, todos, 2)
	assert.Equal(t, "Todo 1", todos[0].Title)
	assert.Equal(t, "Todo 2", todos[1].Title)
	
	mockService.AssertExpectations(t)
}