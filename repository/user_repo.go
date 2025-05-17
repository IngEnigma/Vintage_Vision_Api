package repository

import (
	"context"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"
	"vintage-vision-api/repository/entity"
	"vintage-vision-api/repository/mapper"

	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) domain.UserRepository {
	return &UserRepo{DB: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	entityUser := mapper.FromDomainUser(user)

	err := r.DB.WithContext(ctx).Create(entityUser).Error
	if err != nil {
		utils.Logger.Errorf("%s (%s): %v", constants.ErrMsgCreateUser, user.Email, err)
		return err
	}

	user.ID = entityUser.ID

	utils.Logger.Infof("%s: %s", constants.MsgUserCreatedSuccessfully, user.Email)
	return nil
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var entityUser entity.User
	err := r.DB.WithContext(ctx).
		Preload("Profiles.WatchHistory").
		Where("email = ?", email).
		First(&entityUser).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warnf("%s: %s", constants.ErrMsgUserNotFoundByEmail, email)
		} else {
			utils.Logger.Errorf("%s: %s. Error: %v", constants.ErrMsgFindUserByEmail, email, err)
		}
		return nil, err
	}

	domainUser := mapper.ToDomainUser(&entityUser)

	utils.Logger.Infof("%s: %s", constants.MsgUserFoundByEmail, email)
	return domainUser, nil
}
