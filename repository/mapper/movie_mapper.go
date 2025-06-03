package mapper

import (
	"strconv"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/response"
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

func ToMoviePreview(m domain.Movie) response.MoviePreviewResponse {
	return response.MoviePreviewResponse{
		ID:       strconv.FormatUint(uint64(m.ID), 10),
		ImageURL: m.ImageURL,
	}
}

func ToMoviePreviewList(movies []domain.Movie) []response.MoviePreviewResponse {
	previews := make([]response.MoviePreviewResponse, len(movies))
	for i, m := range movies {
		previews[i] = ToMoviePreview(m)
	}
	return previews
}

func ToMovieDetailResponse(m domain.Movie) response.MovieDetailResponse {
	return response.MovieDetailResponse{
		ID:          strconv.FormatUint(uint64(m.ID), 10),
		Description: m.Description,
		Genre:       m.Genre,
		Year:        m.Year,
		ImageURL:    m.ImageURL,
	}
}

func ToMoviePlayerResponse(m *domain.Movie) response.MoviePlayerResponse {
	return response.MoviePlayerResponse{
		ID:        strconv.FormatUint(uint64(m.ID), 10),
		Title:     m.Title,
		StreamURL: m.StreamURL,
	}
}

func ToMovieTitleResponse(m *domain.Movie) response.MovieTitleResponse {
	return response.MovieTitleResponse{
		ID:    strconv.FormatUint(uint64(m.ID), 10),
		Title: m.Title,
	}
}
