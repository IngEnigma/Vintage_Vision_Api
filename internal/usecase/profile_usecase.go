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

	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = *req.AvatarURL
	}

	if err := u.Repo.Update(ctx, profile); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgUpdateProfile, err)
		return err
	}

	utils.Logger.Infof("%s (%d), perfil ID (%d)", constants.MsgProfileUpdatedSuccessfully, userID, profileID)
	return nil
}
