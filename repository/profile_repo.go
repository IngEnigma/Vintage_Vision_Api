package repository

import (
	"vintage-vision-api/internal/domain"

	"gorm.io/gorm"
)

type ProfileRepo struct {
	DB *gorm.DB
}

func NewProfileRepo(db *gorm.DB) domain.ProfileRepository {
	return &ProfileRepo{DB: db}
}

func (r *ProfileRepo) Create(profile *domain.Profile) error {
	return r.DB.Create(profile).Error
}

func (r *ProfileRepo) FindByUser(userID uint) ([]domain.Profile, error) {
	var profiles []domain.Profile
	err := r.DB.Where("user_id = ?", userID).Find(&profiles).Error
	return profiles, err
}

func (r *ProfileRepo) DeleteByID(profileID uint, userID uint) error {
	result := r.DB.Where("id = ? AND user_id = ?", profileID, userID).Delete(&domain.Profile{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *ProfileRepo) Update(profile *domain.Profile) error {
	return r.DB.Save(profile).Error
}

func (r *ProfileRepo) FindByIDAndUser(profileID, userID uint) (*domain.Profile, error) {
	var profile domain.Profile
	err := r.DB.Where("id = ? AND user_id = ?", profileID, userID).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
