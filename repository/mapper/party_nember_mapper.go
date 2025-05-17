package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"

	"gorm.io/gorm"
)

func ToDomainPartyMember(pm *entity.PartyMember) *domain.PartyMember {
	return &domain.PartyMember{
		ID:        pm.ID,
		PartyID:   pm.PartyID,
		ProfileID: pm.ProfileID,
		JoinedAt:  pm.JoinedAt,
	}
}

func FromDomainPartyMember(pm *domain.PartyMember) *entity.PartyMember {
	return &entity.PartyMember{
		Model:     gorm.Model{ID: pm.ID},
		PartyID:   pm.PartyID,
		ProfileID: pm.ProfileID,
		JoinedAt:  pm.JoinedAt,
	}
}
