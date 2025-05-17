package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"

	"gorm.io/gorm"
)

func ToDomainParty(p *entity.Party) *domain.Party {
	members := make([]domain.PartyMember, len(p.Members))
	for i, m := range p.Members {
		members[i] = *ToDomainPartyMember(&m)
	}

	return &domain.Party{
		ID:        p.ID,
		HostID:    p.HostID,
		MovieID:   p.MovieID,
		PartyCode: p.PartyCode,
		CreatedAt: p.CreatedAt,
		ExpiresAt: p.ExpiresAt,
		Movie:     *ToDomainMovie(&p.Movie),
		Host:      *ToDomainUser(&p.Host),
		Members:   members,
	}
}

func FromDomainParty(p *domain.Party) *entity.Party {
	members := make([]entity.PartyMember, len(p.Members))
	for i, m := range p.Members {
		members[i] = *FromDomainPartyMember(&m)
	}

	return &entity.Party{
		Model:     gorm.Model{ID: p.ID},
		HostID:    p.HostID,
		MovieID:   p.MovieID,
		PartyCode: p.PartyCode,
		ExpiresAt: p.ExpiresAt,
		Movie:     *FromDomainMovie(&p.Movie),
		Host:      *FromDomainUser(&p.Host),
		Members:   members,
	}
}
