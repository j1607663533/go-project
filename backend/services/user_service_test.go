package services

import (
	"testing"

	"gin-backend/models"
	"gin-backend/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository 模拟用户仓储
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// 确保 MockUserRepository 实现了 UserRepository 接口
var _ repositories.UserRepository = (*MockUserRepository)(nil)

func TestGetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUsers := []models.User{
		{ID: 1, Username: "user1", Email: "user1@example.com"},
		{ID: 2, Username: "user2", Email: "user2@example.com"},
	}

	mockRepo.On("FindAll").Return(expectedUsers, nil)

	users, err := service.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := &models.UserCreateRequest{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
		Nickname: "New User",
	}

	mockRepo.On("FindByUsername", "newuser").Return(nil, nil)
	mockRepo.On("FindByEmail", "newuser@example.com").Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	user, err := service.CreateUser(req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_UsernameExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	req := &models.UserCreateRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
	}

	existingUser := &models.User{ID: 1, Username: "existinguser"}
	mockRepo.On("FindByUsername", "existinguser").Return(existingUser, nil)

	user, err := service.CreateUser(req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "用户名已存在", err.Error())
	mockRepo.AssertExpectations(t)
}
