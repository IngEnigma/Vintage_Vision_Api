package usecase

import (
	"errors"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/utils"
)

type UserUsecase struct {
	Repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) *UserUsecase {
	return &UserUsecase{Repo: r}
}

func (u *UserUsecase) Register(req request.RegisterRequest) error {
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashed,
	}

	return u.Repo.Create(user)
}

func (u *UserUsecase) Login(req request.LoginRequest) (string, error) {
	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return "", errors.New("credenciales inválidas")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("credenciales inválidas")
	}

	return utils.GenerateJWT(user.ID, user.IsAdmin)
}
