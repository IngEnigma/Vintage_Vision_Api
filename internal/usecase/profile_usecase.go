package usecase

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
)

type ProfileUsecase struct {
	Repo domain.ProfileRepository
}

func NewProfileUsecase(r domain.ProfileRepository) *ProfileUsecase {
	return &ProfileUsecase{Repo: r}
}

func (u *ProfileUsecase) Create(userID uint, req request.CreateProfileRequest) error {
	profile := &domain.Profile{
		Name:      req.Name,
		AvatarURL: req.AvatarUrl,
		UserID:    userID,
	}
	return u.Repo.Create(profile)
}

func (u *ProfileUsecase) GetAll(userID uint) ([]domain.Profile, error) {
	return u.Repo.FindByUser(userID)
}

func (u *ProfileUsecase) Delete(profileID, userID uint) error {
	return u.Repo.DeleteByID(profileID, userID)
}

func (u *ProfileUsecase) Update(profileID, userID uint, req request.UpdateProfileRequest) error {
	profile, err := u.Repo.FindByIDAndUser(profileID, userID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = *req.AvatarURL
	}

	return u.Repo.Update(profile)
}
