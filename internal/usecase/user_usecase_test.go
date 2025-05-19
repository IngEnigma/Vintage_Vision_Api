package usecase_test

import (
	"context"
	"errors"
	"testing"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

type MockTokenGenerator struct {
	mock.Mock
}

func (m *MockTokenGenerator) Generate(userID uint, isAdmin bool) (string, error) {
	args := m.Called(userID, isAdmin)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	user := args.Get(0)
	if user == nil {
		return nil, args.Error(1)
	}
	return user.(*domain.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	usecase := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.RegisterRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := usecase.Register(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegister_InvalidInput(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	usecase := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.RegisterRequest{
		Email:    "bad-email",
		Password: "123",
	}

	err := usecase.Register(context.Background(), req)

	assert.EqualError(t, err, constants.ErrMsgInvalidRegisterData)
	mockRepo.AssertNotCalled(t, "FindByEmail")
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	usecase := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.RegisterRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).
		Return(&domain.User{Email: req.Email}, nil)

	err := usecase.Register(context.Background(), req)

	assert.EqualError(t, err, constants.ErrMsgEmailAlreadyRegistered)
	mockRepo.AssertExpectations(t)
}

func TestRegister_HashingError(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	usecase := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.RegisterRequest{
		Email:    "test@example.com",
		Password: string(make([]byte, 1<<20)),
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	err := usecase.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bcrypt")
}

func TestRegister_CreateError(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	usecase := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.RegisterRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.User")).
		Return(errors.New("db error"))

	err := usecase.Register(context.Background(), req)

	assert.EqualError(t, err, "db error")
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	userUC := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.LoginRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	hashed, _ := utils.HashPassword(req.Password)

	mockRepo.On("FindByEmail", mock.Anything, req.Email).
		Return(&domain.User{ID: 1, Email: req.Email, Password: hashed, IsAdmin: false}, nil)

	mockTokenGen.On("Generate", uint(1), false).Return("fake-jwt-token", nil)

	token, err := userUC.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}

func TestLogin_InvalidInput(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	userUC := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.LoginRequest{
		Email:    "bad-email",
		Password: "",
	}

	token, err := userUC.Login(context.Background(), req)

	assert.EqualError(t, err, constants.ErrMsgInvalidLoginData)
	assert.Empty(t, token)
	mockRepo.AssertNotCalled(t, "FindByEmail")
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	userUC := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.LoginRequest{
		Email:    "notfound@example.com",
		Password: "securepassword",
	}

	mockRepo.On("FindByEmail", mock.Anything, req.Email).Return(nil, nil)

	token, err := userUC.Login(context.Background(), req)

	assert.EqualError(t, err, constants.ErrMsgInvalidCredentials)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)
	userUC := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	hashed, _ := utils.HashPassword("correctpassword")

	mockRepo.On("FindByEmail", mock.Anything, req.Email).
		Return(&domain.User{Email: req.Email, Password: hashed}, nil)

	token, err := userUC.Login(context.Background(), req)

	assert.EqualError(t, err, constants.ErrMsgInvalidCredentials)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_GenerateTokenError(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockTokenGen := new(MockTokenGenerator)

	userUC := usecase.NewUserUsecase(mockRepo, mockTokenGen)

	req := request.LoginRequest{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	hashed, _ := utils.HashPassword(req.Password)

	mockRepo.On("FindByEmail", mock.Anything, req.Email).
		Return(&domain.User{ID: 1, Email: req.Email, Password: hashed, IsAdmin: false}, nil)

	mockTokenGen.On("Generate", uint(1), false).
		Return("", errors.New("token error"))

	token, err := userUC.Login(context.Background(), req)

	assert.EqualError(t, err, constants.ErrMsgGeneratingToken)
	assert.Empty(t, token)

	mockRepo.AssertExpectations(t)
	mockTokenGen.AssertExpectations(t)
}
