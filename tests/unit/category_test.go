package unit

import (
	"testing"
	"time"

	"mcp-godo/pkg/todo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository is a mock implementation of CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(category todo.Category) (todo.Category, error) {
	args := m.Called(category)
	return args.Get(0).(todo.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindAll() ([]todo.Category, error) {
	args := m.Called()
	return args.Get(0).([]todo.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByID(id int64) (todo.Category, error) {
	args := m.Called(id)
	return args.Get(0).(todo.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByName(name string) (todo.Category, error) {
	args := m.Called(name)
	return args.Get(0).(todo.Category), args.Error(1)
}

func (m *MockCategoryRepository) Update(category todo.Category) (todo.Category, error) {
	args := m.Called(category)
	return args.Get(0).(todo.Category), args.Error(1)
}

func (m *MockCategoryRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCategoryRepository) FindTodosByCategory(categoryID int64) ([]todo.TodoItem, error) {
	args := m.Called(categoryID)
	return args.Get(0).([]todo.TodoItem), args.Error(1)
}

func (m *MockCategoryRepository) FindUncategorizedTodos() ([]todo.TodoItem, error) {
	args := m.Called()
	return args.Get(0).([]todo.TodoItem), args.Error(1)
}

func TestCreateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	testCategory := todo.Category{
		Name:        "Work Tasks",
		Description: stringPtr("Professional tasks"),
		Color:       stringPtr("#3498db"),
	}

	expectedCategory := todo.Category{
		ID:          1,
		Name:        "Work Tasks",
		Description: stringPtr("Professional tasks"),
		Color:       stringPtr("#3498db"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock the FindByName call to return not found
	mockRepo.On("FindByName", "Work Tasks").Return(todo.Category{}, assert.AnError)
	// Mock the Create call
	mockRepo.On("Create", testCategory).Return(expectedCategory, nil)

	result, err := service.CreateCategory("Work Tasks", stringPtr("Professional tasks"), stringPtr("#3498db"))

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory.ID, result.ID)
	assert.Equal(t, expectedCategory.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateCategory_EmptyName(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	_, err := service.CreateCategory("", nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category name cannot be empty")
	mockRepo.AssertExpectations(t)
}

func TestCreateCategory_DuplicateName(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	existingCategory := todo.Category{
		ID:   1,
		Name: "Work Tasks",
	}

	// Mock the FindByName call to return existing category
	mockRepo.On("FindByName", "Work Tasks").Return(existingCategory, nil)

	_, err := service.CreateCategory("Work Tasks", nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category with name 'Work Tasks' already exists")
	mockRepo.AssertExpectations(t)
}

func TestCreateCategory_InvalidColor(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	invalidColor := "invalid-color"

	_, err := service.CreateCategory("Work Tasks", nil, &invalidColor)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid hex color format")
	mockRepo.AssertExpectations(t)
}

func TestGetAllCategories_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	expectedCategories := []todo.Category{
		{ID: 1, Name: "Work Tasks"},
		{ID: 2, Name: "Personal Tasks"},
	}

	mockRepo.On("FindAll").Return(expectedCategories, nil)

	result, err := service.GetAllCategories()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

func TestGetCategoryByID_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	expectedCategory := todo.Category{
		ID:   1,
		Name: "Work Tasks",
	}

	mockRepo.On("FindByID", int64(1)).Return(expectedCategory, nil)

	result, err := service.GetCategoryByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory.ID, result.ID)
	assert.Equal(t, expectedCategory.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	existingCategory := todo.Category{
		ID:   1,
		Name: "Old Name",
	}

	updatedCategory := todo.Category{
		ID:   1,
		Name: "New Name",
	}

	newName := "New Name"

	// Mock the FindByID call
	mockRepo.On("FindByID", int64(1)).Return(existingCategory, nil)
	// Mock the Update call
	mockRepo.On("Update", updatedCategory).Return(updatedCategory, nil)

	result, err := service.UpdateCategory(1, &newName, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, "New Name", result.Name)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	mockRepo.On("Delete", int64(1)).Return(nil)

	err := service.DeleteCategory(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTodosByCategory_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	expectedTodos := []todo.TodoItem{
		{ID: "1", Title: "Task 1", CategoryID: int64Ptr(1)},
		{ID: "2", Title: "Task 2", CategoryID: int64Ptr(1)},
	}

	mockRepo.On("FindTodosByCategory", int64(1)).Return(expectedTodos, nil)

	result, err := service.GetTodosByCategory(1)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

func TestGetUncategorizedTodos_Success(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	service := todo.NewCategoryService(mockRepo)

	expectedTodos := []todo.TodoItem{
		{ID: "1", Title: "Task 1", CategoryID: nil},
		{ID: "2", Title: "Task 2", CategoryID: nil},
	}

	mockRepo.On("FindUncategorizedTodos").Return(expectedTodos, nil)

	result, err := service.GetUncategorizedTodos()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	mockRepo.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}