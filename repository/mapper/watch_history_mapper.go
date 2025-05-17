package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"
)

func ToDomainWatchHistory(e *entity.WatchHistory) *domain.WatchHistory {
	return &domain.WatchHistory{
		ID:        e.ID,
		ProfileID: e.ProfileID,
		MovieID:   e.MovieID,
		WatchedAt: e.WatchedAt,
		Progress:  e.Progress,
	}
}

func FromDomainWatchHistory(d *domain.WatchHistory) *entity.WatchHistory {
	return &entity.WatchHistory{
		ID:        d.ID,
		ProfileID: d.ProfileID,
		MovieID:   d.MovieID,
		WatchedAt: d.WatchedAt,
		Progress:  d.Progress,
	}
}
