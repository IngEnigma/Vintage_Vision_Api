package usecase_test

/*
import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"
)

type MockUserRepo struct {
	mock.Mock
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

func TestRegister(t *testing.T) {
	ctx := context.Background()

	t.Run("email ya registrado", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		uc := usecase.NewUserUsecase(mockRepo)

		req := request.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", ctx, req.Email).Return(&domain.User{}, nil)

		err := uc.Register(ctx, req)

		assert.EqualError(t, err, constants.ErrMsgEmailAlreadyRegistered)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fallo al hashear la contraseña", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		uc := usecase.NewUserUsecase(mockRepo)

		req := request.RegisterRequest{
			Email:    "failhash@example.com",
			Password: string(make([]byte, 0)),
		}

		mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, nil)

		err := uc.Register(ctx, req)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error al crear el usuario", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		uc := usecase.NewUserUsecase(mockRepo)

		req := request.RegisterRequest{
			Email:    "new@example.com",
			Password: "securepassword",
		}

		mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, nil)

		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(errors.New("db error"))

		err := uc.Register(ctx, req)

		assert.EqualError(t, err, "db error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("registro exitoso", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		uc := usecase.NewUserUsecase(mockRepo)

		req := request.RegisterRequest{
			Email:    "success@example.com",
			Password: "validpassword",
		}

		mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

		err := uc.Register(ctx, req)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("usuario no encontrado", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		uc := usecase.NewUserUsecase(mockRepo)

		req := request.LoginRequest{
			Email:    "notfound@example.com",
			Password: "password123",
		}

		mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.New("not found"))

		token, err := uc.Login(ctx, req)

		assert.Empty(t, token)
		assert.EqualError(t, err, constants.ErrMsgInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contraseña incorrecta", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		uc := usecase.NewUserUsecase(mockRepo)

		req := request.LoginRequest{
			Email:    "user@example.com",
			Password: "wrongpass",
		}

		hashed, _ := utils.HashPassword("correctpass")

		mockRepo.On("FindByEmail", ctx, req.Email).Return(&domain.User{
			Email:    req.Email,
			Password: hashed,
		}, nil)

		token, err := uc.Login(ctx, req)

		assert.Empty(t, token)
		assert.EqualError(t, err, constants.ErrMsgInvalidCredentials)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error al generar token", func(t *testing.T) {
		t.Skip("No se puede simular fallo en GenerateJWT sin modificar el código de producción.")
	})

	t.Run("login exitoso", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		uc := usecase.NewUserUsecase(mockRepo)

		req := request.LoginRequest{
			Email:    "user@example.com",
			Password: "mypassword",
		}

		hashed, _ := utils.HashPassword(req.Password)

		mockRepo.On("FindByEmail", ctx, req.Email).Return(&domain.User{
			ID:       1,
			Email:    req.Email,
			Password: hashed,
			IsAdmin:  false,
		}, nil)

		token, err := uc.Login(ctx, req)

		assert.NotEmpty(t, token)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
*/
