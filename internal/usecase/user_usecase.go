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
	existingUser, _ := u.Repo.FindByEmail(req.Email)
	if existingUser != nil {
		utils.Logger.Warnf("Registro fallido: el correo ya está registrado - %s", req.Email)
		return errors.New("el correo ya está registrado")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Logger.Errorf("Error al hashear la contraseña para %s: %v", req.Email, err)
		return err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashed,
	}

	if err := u.Repo.Create(user); err != nil {
		utils.Logger.Errorf("Error al crear el usuario en la base de datos (%s): %v", req.Email, err)
		return err
	}

	utils.Logger.Infof("Usuario creado exitosamente en base de datos: %s", req.Email)
	return nil
}

func (u *UserUsecase) Login(req request.LoginRequest) (string, error) {
	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		utils.Logger.Warnf("Inicio de sesión fallido - usuario no encontrado: %s", req.Email)
		return "", errors.New("credenciales inválidas")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.Logger.Warnf("Inicio de sesión fallido - contraseña incorrecta: %s", req.Email)
		return "", errors.New("credenciales inválidas")
	}

	token, err := utils.GenerateJWT(user.ID, user.IsAdmin)
	if err != nil {
		utils.Logger.Errorf("Error generando token JWT para %s: %v", req.Email, err)
		return "", errors.New("error generando token")
	}

	utils.Logger.Infof("Inicio de sesión exitoso: %s", req.Email)
	return token, nil
}
