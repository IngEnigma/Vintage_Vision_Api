package repository

import (
	"context"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) domain.UserRepository {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	err := r.DB.WithContext(ctx).Create(user).Error
	if err != nil {
		utils.Logger.Errorf("%s (%s): %v", constants.ErrMsgCreateUser, user.Email, err)
	} else {
		utils.Logger.Infof("%s: %s", constants.MsgUserCreatedSuccessfully, user.Email)
	}
	return err
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warnf("%s: %s", constants.ErrMsgUserNotFoundByEmail, email)
		} else {
			utils.Logger.Errorf("%s: %s. Error: %v", constants.ErrMsgFindUserByEmail, email, err)
		}
		return nil, err
	}

	utils.Logger.Infof("%s: %s", constants.MsgUserFoundByEmail, email)
	return &user, nil
}
