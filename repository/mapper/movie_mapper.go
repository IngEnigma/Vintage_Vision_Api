package mapper

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/repository/entity"
)

func ToDomainMovie(e *entity.Movie) *domain.Movie {
	watchHistories := make([]domain.WatchHistory, len(e.WatchHistories))
	for i, wh := range e.WatchHistories {
		watchHistories[i] = *ToDomainWatchHistory(&wh)
	}

	return &domain.Movie{
		ID:             e.ID,
		Title:          e.Title,
		Description:    e.Description,
		Year:           e.Year,
		Genre:          e.Genre,
		ImageURL:       e.ImageURL,
		StreamURL:      e.StreamURL,
		Duration:       e.Duration,
		WatchHistories: watchHistories,
	}
}

func FromDomainMovie(d *domain.Movie) *entity.Movie {
	watchHistories := make([]entity.WatchHistory, len(d.WatchHistories))
	for i, wh := range d.WatchHistories {
		watchHistories[i] = *FromDomainWatchHistory(&wh)
	}

	return &entity.Movie{
		Title:          d.Title,
		Description:    d.Description,
		Year:           d.Year,
		Genre:          d.Genre,
		ImageURL:       d.ImageURL,
		StreamURL:      d.StreamURL,
		Duration:       d.Duration,
		WatchHistories: watchHistories,
	}
}
