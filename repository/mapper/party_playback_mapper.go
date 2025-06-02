package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"

	"gorm.io/gorm"
)

func ToDomainPartyPlayback(pp *entity.PartyPlayback) *domain.PartyPlayback {
	return &domain.PartyPlayback{
		ID:        pp.ID,
		PartyID:   pp.PartyID,
		IsPlaying: pp.IsPlaying,
		CurrentAt: pp.CurrentAt,
		UpdatedAt: pp.UpdatedAt,
	}
}

func FromDomainPartyPlayback(pp *domain.PartyPlayback) *entity.PartyPlayback {
	return &entity.PartyPlayback{
		Model:     gorm.Model{ID: pp.ID},
		PartyID:   pp.PartyID,
		IsPlaying: pp.IsPlaying,
		CurrentAt: pp.CurrentAt,
		UpdatedAt: pp.UpdatedAt,
	}
}
