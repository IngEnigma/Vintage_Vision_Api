package usecase

import (
	"context"
	"errors"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/utils"
)

type ProfileUsecase struct {
	Repo domain.ProfileRepository
}

func NewProfileUsecase(r domain.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{Repo: r}
}

func (u *ProfileUsecase) Create(ctx context.Context, userID uint, req request.CreateProfileRequest) error {
	if !utils.IsValidProfileName(req.Name) {
		return errors.New(constants.ErrMsgInvalidProfileName)
	}

	if !utils.IsValidURL(req.AvatarUrl) {
		return errors.New(constants.ErrMsgInvalidAvatarURL)
	}

	existingProfiles, err := u.Repo.FindByUser(ctx, userID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetProfiles, err)
		return err
	}
	if len(existingProfiles) >= 4 {
		return errors.New(constants.ErrMsgMaxProfilesReached)
	}

	for _, p := range existingProfiles {
		if p.Name == req.Name {
			return errors.New(constants.ErrMsgProfileNameTaken)
		}
	}

	profile := &domain.Profile{
		Name:      req.Name,
		AvatarURL: req.AvatarUrl,
		UserID:    userID,
	}

	if err := u.Repo.Create(ctx, profile); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgRegisterProfile, err)
		return err
	}

	utils.Logger.Infof("%s (%d)", constants.MsgProfileCreatedSuccessfully, userID)
	return nil
}

func (u *ProfileUsecase) GetAll(ctx context.Context, userID uint) ([]domain.Profile, error) {
	profiles, err := u.Repo.FindByUser(ctx, userID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetProfiles, err)
		return nil, err
	}

	utils.Logger.Infof("%s (%d)", constants.MsgProfilesRetrievedSuccessfully, userID)
	return profiles, nil
}

func (u *ProfileUsecase) Delete(ctx context.Context, profileID, userID uint) error {
	profile, err := u.Repo.FindByIDAndUser(ctx, profileID, userID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetProfile, err)
		return err
	}
	if profile == nil {
		utils.Logger.Warnf("%s (userID: %d, profileID: %d)", constants.ErrMsgProfileNotFound, userID, profileID)
		return errors.New(constants.ErrMsgProfileNotFound)
	}

	if err := u.Repo.DeleteByID(ctx, profileID, userID); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgDeleteProfile, err)
		return err
	}

	utils.Logger.Infof("%s (%d), perfil ID (%d)", constants.MsgProfileDeletedSuccessfully, userID, profileID)
	return nil
}

func (u *ProfileUsecase) Update(ctx context.Context, profileID, userID uint, req request.UpdateProfileRequest) error {
	profile, err := u.Repo.FindByIDAndUser(ctx, profileID, userID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetProfile, err)
		return err
	}
	if profile == nil {
		utils.Logger.Warnf("%s (userID: %d, profileID: %d)", constants.ErrMsgProfileNotFound, userID, profileID)
		return errors.New(constants.ErrMsgProfileNotFound)
	}

	existingProfiles, err := u.Repo.FindByUser(ctx, userID)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetProfiles, err)
		return err
	}

	if req.Name != nil {
		if !utils.IsValidProfileName(*req.Name) {
			return errors.New(constants.ErrMsgInvalidProfileName)
		}

		for _, p := range existingProfiles {
			if p.Name == *req.Name && p.ID != profileID {
				return errors.New(constants.ErrMsgProfileNameTaken)
			}
		}

		profile.Name = *req.Name
	}

	if req.AvatarURL != nil {
		if !utils.IsValidURL(*req.AvatarURL) {
			return errors.New(constants.ErrMsgInvalidAvatarURL)
		}
		profile.AvatarURL = *req.AvatarURL
	}

	if err := u.Repo.Update(ctx, profile); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgUpdateProfile, err)
		return err
	}

	utils.Logger.Infof("%s (%d), perfil ID (%d)", constants.MsgProfileUpdatedSuccessfully, userID, profileID)
	return nil
}
