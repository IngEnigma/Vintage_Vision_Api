package repository

import (
	"vintage-vision-api/internal/domain"

	"gorm.io/gorm"
)

type MovieRepo struct {
	DB *gorm.DB
}

func NewMovieRepo(db *gorm.DB) domain.MovieRepository {
	return &MovieRepo{DB: db}
}

func (r *MovieRepo) GetAll() ([]domain.Movie, error) {
	var movies []domain.Movie
	err := r.DB.Find(&movies).Error
	return movies, err
}

func (r *MovieRepo) Create(movie *domain.Movie) error {
	return r.DB.Create(movie).Error
}

func (r *MovieRepo) GetByID(id uint) (*domain.Movie, error) {
	var movie domain.Movie
	err := r.DB.First(&movie, id).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepo) Update(movie *domain.Movie) error {
	return r.DB.Save(movie).Error
}

func (r *MovieRepo) Delete(id uint) error {
	return r.DB.Delete(&domain.Movie{}, id).Error
}
