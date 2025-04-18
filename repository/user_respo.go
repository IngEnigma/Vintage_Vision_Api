package repository

import (
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

func (r *UserRepo) Create(user *domain.User) error {
	err := r.DB.Create(user).Error
	if err != nil {
		utils.Logger.Errorf("Error al crear el usuario en la base de datos (%s): %v", user.Email, err)
	} else {
		utils.Logger.Infof("Usuario creado en la base de datos: %s", user.Email)
	}
	return err
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warnf("Usuario no encontrado con el email: %s", email)
		} else {
			utils.Logger.Errorf("Error al buscar usuario con el email: %s. Error: %v", email, err)
		}
		return nil, err
	}
	utils.Logger.Infof("Usuario encontrado por email: %s", email)
	return &user, nil
}
