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

var ErrProfileNotFound = gorm.ErrRecordNotFound

type ProfileRepo struct {
	DB *gorm.DB
}

func NewProfileRepo(db *gorm.DB) domain.ProfileRepository {
	return &ProfileRepo{DB: db}
}

func (r *ProfileRepo) Create(ctx context.Context, profile *domain.Profile) error {
	entityProfile := mapper.FromDomainProfile(profile)
	err := r.DB.WithContext(ctx).Create(entityProfile).Error
	if err != nil {
		utils.Logger.Errorf("%s: userID=%d, error=%v", constants.ErrMsgRegisterProfile, profile.UserID, err)
		return err
	}

	profile.ID = entityProfile.ID

	utils.Logger.Infof("%s: userID=%d", constants.MsgProfileCreatedSuccessfully, profile.UserID)
	return nil
}

func (r *ProfileRepo) FindByUser(ctx context.Context, userID uint) ([]domain.Profile, error) {
	var entityProfiles []entity.Profile
	err := r.DB.WithContext(ctx).Preload("WatchHistory").Where("user_id = ?", userID).Find(&entityProfiles).Error
	if err != nil {
		utils.Logger.Errorf("%s: userID=%d, error=%v", constants.ErrMsgGetProfiles, userID, err)
		return nil, err
	}

	domainProfiles := make([]domain.Profile, len(entityProfiles))
	for i, ep := range entityProfiles {
		domainProfiles[i] = *mapper.ToDomainProfile(&ep)
	}

	utils.Logger.Infof("%s: userID=%d", constants.MsgProfilesRetrievedSuccessfully, userID)
	return domainProfiles, nil
}

func (r *ProfileRepo) DeleteByID(ctx context.Context, profileID, userID uint) error {
	result := r.DB.WithContext(ctx).Where("id = ? AND user_id = ?", profileID, userID).Delete(&entity.Profile{})
	if result.Error != nil {
		utils.Logger.Errorf("%s: profileID=%d, userID=%d, error=%v", constants.ErrMsgDeleteProfile, profileID, userID, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		utils.Logger.Warnf("%s: profileID=%d, userID=%d", constants.ErrMsgProfileNotFound, profileID, userID)
		return ErrProfileNotFound
	}

	utils.Logger.Infof("%s: profileID=%d, userID=%d", constants.MsgProfileDeletedSuccessfully, profileID, userID)
	return nil
}

func (r *ProfileRepo) Update(ctx context.Context, profile *domain.Profile) error {
	entityProfile := mapper.FromDomainProfile(profile)
	result := r.DB.WithContext(ctx).Save(entityProfile)
	if result.Error != nil {
		utils.Logger.Errorf("%s: profileID=%d, error=%v", constants.ErrMsgUpdateProfile, profile.ID, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		utils.Logger.Warnf("%s: profileID=%d", constants.ErrMsgProfileNotFound, profile.ID)
		return ErrProfileNotFound
	}

	utils.Logger.Infof("%s: profileID=%d", constants.MsgProfileUpdatedSuccessfully, profile.ID)
	return nil
}

func (r *ProfileRepo) FindByIDAndUser(ctx context.Context, profileID, userID uint) (*domain.Profile, error) {
	var entityProfile entity.Profile
	err := r.DB.WithContext(ctx).Preload("WatchHistory").
		Where("id = ? AND user_id = ?", profileID, userID).
		First(&entityProfile).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warnf("%s: profileID=%d, userID=%d", constants.ErrMsgProfileNotFound, profileID, userID)
			return nil, ErrProfileNotFound
		}
		utils.Logger.Errorf("%s: profileID=%d, userID=%d, error=%v", constants.ErrMsgGetProfile, profileID, userID, err)
		return nil, err
	}

	domainProfile := mapper.ToDomainProfile(&entityProfile)
	utils.Logger.Infof("%s: profileID=%d, userID=%d", constants.MsgProfileRetrievedSuccessfully, profileID, userID)
	return domainProfile, nil
}
