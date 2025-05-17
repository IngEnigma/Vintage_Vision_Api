package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"
)

func ToDomainUser(u *entity.User) *domain.User {
	profiles := make([]domain.Profile, len(u.Profiles))
	for i, p := range u.Profiles {
		profiles[i] = *ToDomainProfile(&p)
	}

	return &domain.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
		Profiles: profiles,
	}
}

func FromDomainUser(u *domain.User) *entity.User {
	profiles := make([]entity.Profile, len(u.Profiles))
	for i, p := range u.Profiles {
		profiles[i] = *FromDomainProfile(&p)
	}

	return &entity.User{
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
		Profiles: profiles,
	}
}
