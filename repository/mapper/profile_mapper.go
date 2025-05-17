package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"

	"gorm.io/gorm"
)

func ToDomainProfile(p *entity.Profile) *domain.Profile {
	watchHistory := make([]domain.WatchHistory, len(p.WatchHistory))
	for i, wh := range p.WatchHistory {
		watchHistory[i] = *ToDomainWatchHistory(&wh)
	}

	return &domain.Profile{
		ID:           p.ID,
		Name:         p.Name,
		AvatarURL:    p.AvatarURL,
		UserID:       p.UserID,
		WatchHistory: watchHistory,
	}
}

func FromDomainProfile(p *domain.Profile) *entity.Profile {
	watchHistory := make([]entity.WatchHistory, len(p.WatchHistory))
	for i, wh := range p.WatchHistory {
		watchHistory[i] = *FromDomainWatchHistory(&wh)
	}

	return &entity.Profile{
		Model:        gorm.Model{ID: p.ID},
		Name:         p.Name,
		AvatarURL:    p.AvatarURL,
		UserID:       p.UserID,
		WatchHistory: watchHistory,
	}
}
