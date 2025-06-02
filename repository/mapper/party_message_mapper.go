package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"

	"gorm.io/gorm"
)

func ToDomainPartyMessage(pm *entity.PartyMessage) *domain.PartyMessage {
	return &domain.PartyMessage{
		ID:       pm.ID,
		PartyID:  pm.PartyID,
		SenderID: pm.SenderID,
		Message:  pm.Message,
		SentAt:   pm.SentAt,
		Sender:   *ToDomainProfile(&pm.Sender),
	}
}

func FromDomainPartyMessage(pm *domain.PartyMessage) *entity.PartyMessage {
	return &entity.PartyMessage{
		Model:    gorm.Model{ID: pm.ID},
		PartyID:  pm.PartyID,
		SenderID: pm.SenderID,
		Message:  pm.Message,
		SentAt:   pm.SentAt,
		Sender:   *FromDomainProfile(&pm.Sender),
	}
}
