package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"

	"gorm.io/gorm"
)

func ToDomainPartyMediaStatus(pms *entity.PartyMediaStatus) *domain.PartyMediaStatus {
	return &domain.PartyMediaStatus{
		ID:        pms.ID,
		PartyID:   pms.PartyID,
		ProfileID: pms.ProfileID,
		CameraOn:  pms.CameraOn,
		MicOn:     pms.MicOn,
		UpdatedAt: pms.UpdatedAt,
	}
}

func FromDomainPartyMediaStatus(pms *domain.PartyMediaStatus) *entity.PartyMediaStatus {
	return &entity.PartyMediaStatus{
		Model:     gorm.Model{ID: pms.ID},
		PartyID:   pms.PartyID,
		ProfileID: pms.ProfileID,
		CameraOn:  pms.CameraOn,
		MicOn:     pms.MicOn,
		UpdatedAt: pms.UpdatedAt,
	}
}
