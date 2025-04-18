package usecase

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
)

type MovieUsecase struct {
	Repo domain.MovieRepository
}

func NewMovieUsecase(r domain.MovieRepository) *MovieUsecase {
	return &MovieUsecase{Repo: r}
}

func (u *MovieUsecase) GetAll() ([]domain.Movie, error) {
	return u.Repo.GetAll()
}

func (u *MovieUsecase) Create(req request.CreateMovieRequest) error {
	movie := domain.Movie{
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		ImageURL:    req.ImageURL,
		StreamURL:   req.StreamURL,
		Genre:       req.Genre,
		Duration:    req.Duration,
	}
	return u.Repo.Create(&movie)
}

func (u *MovieUsecase) Update(id uint, req request.UpdateMovieRequest) error {
	movie, err := u.Repo.GetByID(id)
	if err != nil {
		return err
	}

	movie.Title = req.Title
	movie.Description = req.Description
	movie.Year = req.Year
	movie.ImageURL = req.ImageURL
	movie.StreamURL = req.StreamURL
	movie.Genre = req.Genre
	movie.Duration = req.Duration

	return u.Repo.Update(movie)
}

func (u *MovieUsecase) Delete(id uint) error {
	return u.Repo.Delete(id)
}
