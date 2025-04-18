package repository

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"

	"gorm.io/gorm"
)

type ProfileRepo struct {
	DB *gorm.DB
}

func NewProfileRepo(db *gorm.DB) domain.ProfileRepository {
	return &ProfileRepo{DB: db}
}

func (r *ProfileRepo) Create(profile *domain.Profile) error {
	err := r.DB.Create(profile).Error
	if err != nil {
		utils.Logger.Errorf("Error al crear el perfil para el usuario ID (%d): %v", profile.UserID, err)
	} else {
		utils.Logger.Infof("Perfil creado para el usuario ID (%d)", profile.UserID)
	}
	return err
}

func (r *ProfileRepo) FindByUser(userID uint) ([]domain.Profile, error) {
	var profiles []domain.Profile
	err := r.DB.Where("user_id = ?", userID).Find(&profiles).Error
	if err != nil {
		utils.Logger.Errorf("Error al obtener perfiles para el usuario ID (%d): %v", userID, err)
	} else {
		utils.Logger.Infof("Perfiles obtenidos para el usuario ID (%d)", userID)
	}
	return profiles, err
}

func (r *ProfileRepo) DeleteByID(profileID uint, userID uint) error {
	result := r.DB.Where("id = ? AND user_id = ?", profileID, userID).Delete(&domain.Profile{})
	if result.RowsAffected == 0 {
		utils.Logger.Warnf("No se encontró el perfil con ID (%d) para el usuario ID (%d) o no se tiene permiso para eliminarlo", profileID, userID)
		return gorm.ErrRecordNotFound
	}

	utils.Logger.Infof("Perfil con ID (%d) eliminado para el usuario ID (%d)", profileID, userID)
	return result.Error
}

func (r *ProfileRepo) Update(profile *domain.Profile) error {
	err := r.DB.Save(profile).Error
	if err != nil {
		utils.Logger.Errorf("Error al actualizar el perfil con ID (%d): %v", profile.ID, err)
	} else {
		utils.Logger.Infof("Perfil con ID (%d) actualizado", profile.ID)
	}
	return err
}

func (r *ProfileRepo) FindByIDAndUser(profileID, userID uint) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.DB.Where("id = ? AND user_id = ?", profileID, userID).First(&profile).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warnf("Perfil no encontrado con ID (%d) para el usuario ID (%d)", profileID, userID)
		} else {
			utils.Logger.Errorf("Error al obtener el perfil con ID (%d) para el usuario ID (%d): %v", profileID, userID, err)
		}
		return nil, err
	}

	utils.Logger.Infof("Perfil con ID (%d) encontrado para el usuario ID (%d)", profileID, userID)
	return &profile, nil
}
