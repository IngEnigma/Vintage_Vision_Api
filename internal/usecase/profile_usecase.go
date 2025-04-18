package usecase

import (
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

func (u *ProfileUsecase) Create(userID uint, req request.CreateProfileRequest) error {
	profile := &domain.Profile{
		Name:      req.Name,
		AvatarURL: req.AvatarUrl,
		UserID:    userID,
	}
	err := u.Repo.Create(profile)
	if err != nil {
		utils.Logger.Errorf("Error al crear el perfil para el usuario ID (%d): %v", userID, err)
	} else {
		utils.Logger.Infof("Perfil creado con éxito para el usuario ID (%d)", userID)
	}
	return err
}

func (u *ProfileUsecase) GetAll(userID uint) ([]domain.Profile, error) {
	profiles, err := u.Repo.FindByUser(userID)
	if err != nil {
		utils.Logger.Errorf("Error al obtener perfiles para el usuario ID (%d): %v", userID, err)
	} else {
		utils.Logger.Infof("Perfiles obtenidos con éxito para el usuario ID (%d)", userID)
	}
	return profiles, err
}

func (u *ProfileUsecase) Delete(profileID, userID uint) error {
	err := u.Repo.DeleteByID(profileID, userID)
	if err != nil {
		utils.Logger.Errorf("Error al eliminar el perfil con ID (%d) para el usuario ID (%d): %v", profileID, userID, err)
	} else {
		utils.Logger.Infof("Perfil con ID (%d) eliminado con éxito para el usuario ID (%d)", profileID, userID)
	}
	return err
}

func (u *ProfileUsecase) Update(profileID, userID uint, req request.UpdateProfileRequest) error {
	profile, err := u.Repo.FindByIDAndUser(profileID, userID)
	if err != nil {
		utils.Logger.Errorf("Error al obtener el perfil con ID (%d) para el usuario ID (%d): %v", profileID, userID, err)
		return err
	}

	if req.Name != nil {
		profile.Name = *req.Name
	}
	if req.AvatarURL != nil {
		profile.AvatarURL = *req.AvatarURL
	}

	err = u.Repo.Update(profile)
	if err != nil {
		utils.Logger.Errorf("Error al actualizar el perfil con ID (%d) para el usuario ID (%d): %v", profileID, userID, err)
	} else {
		utils.Logger.Infof("Perfil con ID (%d) actualizado con éxito para el usuario ID (%d)", profileID, userID)
	}
	return err
}
