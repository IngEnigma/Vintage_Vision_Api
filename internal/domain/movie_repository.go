package domain

import (
	"context"
	"vintage-vision-api/internal/model/response"
)

type MovieRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]Movie, error)
	GetByID(ctx context.Context, id uint) (*Movie, error)
	Create(ctx context.Context, movie *Movie) error
	Update(ctx context.Context, movie *Movie) error
	UpdateFields(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
	GetMoviesByGenre(ctx context.Context, genre string, page int, limit int) ([]Movie, error)
	GetMoviesByTitle(ctx context.Context, title string, page int, limit int) ([]response.MovieTitleResponse, error)
}
