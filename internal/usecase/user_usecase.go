package usecase

import (
	"context"
	"errors"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/utils"
)

type UserUsecase struct {
	Repo  domain.UserRepository
	Token domain.TokenGenerator
}

func NewUserUsecase(r domain.UserRepository, t domain.TokenGenerator) *UserUsecase {
	return &UserUsecase{Repo: r, Token: t}
}

func (u *UserUsecase) Register(ctx context.Context, req request.RegisterRequest) error {
	if errs := utils.ValidateRegisterInput(req); len(errs) > 0 {
		utils.Logger.Warnf("Errores de validación en registro: %+v", errs)
		return errors.New(constants.ErrMsgInvalidRegisterData)
	}

	existingUser, _ := u.Repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		utils.Logger.Warnf("%s: %s", constants.ErrMsgEmailAlreadyRegistered, req.Email)
		return errors.New(constants.ErrMsgEmailAlreadyRegistered)
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgHashingPassword, err)
		return err
	}

	user := &domain.User{
		Email:    req.Email,
		Password: hashed,
	}

	if err := u.Repo.Create(ctx, user); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCreateUser, err)
		return err
	}

	utils.Logger.Infof("%s: %s", constants.MsgUserCreatedSuccessfully, req.Email)
	return nil
}

func (u *UserUsecase) Login(ctx context.Context, req request.LoginRequest) (string, error) {
	if errs := utils.ValidateLoginInput(req); len(errs) > 0 {
		utils.Logger.Warnf("Errores de validación en login: %+v", errs)
		return "", errors.New(constants.ErrMsgInvalidLoginData)
	}

	user, err := u.Repo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		utils.Logger.Warnf("%s: %s", constants.ErrMsgUserNotFound, req.Email)
		return "", errors.New(constants.ErrMsgInvalidCredentials)
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.Logger.Warnf("%s: %s", constants.ErrMsgInvalidCredentials, req.Email)
		return "", errors.New(constants.ErrMsgInvalidCredentials)
	}

	token, err := u.Token.Generate(user.ID, user.IsAdmin)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGeneratingToken, err)
		return "", errors.New(constants.ErrMsgGeneratingToken)
	}

	utils.Logger.Infof("%s: %s", constants.MsgLoginSuccessful, req.Email)
	return token, nil
}
