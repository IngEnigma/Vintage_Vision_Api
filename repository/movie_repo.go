package repository

import (
	"context"
	"errors"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"

	"gorm.io/gorm"
)

var ErrMovieNotFound = errors.New("movie not found")

type MovieRepo struct {
	DB *gorm.DB
}

func NewMovieRepo(db *gorm.DB) domain.MovieRepository {
	return &MovieRepo{DB: db}
}

func (r *MovieRepo) GetAll(ctx context.Context, limit, offset int) ([]domain.Movie, error) {
	var movies []domain.Movie
	err := r.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&movies).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetMovies, err)
		return nil, err
	}
	utils.Logger.Infof("%s", constants.MsgMoviesRetrievedSuccessfully)
	return movies, nil
}

func (r *MovieRepo) Create(ctx context.Context, movie *domain.Movie) error {
	if err := movie.Validate(); err != nil {
		utils.Logger.Warnf("validation failed: %v", err)
		return err
	}
	err := r.DB.WithContext(ctx).Create(movie).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCreateMovie, err)
		return err
	}
	utils.Logger.Infof("%s: ID %d", constants.MsgMovieCreatedSuccessfully, movie.ID)
	return nil
}

func (r *MovieRepo) Update(ctx context.Context, movie *domain.Movie) error {
	if err := movie.Validate(); err != nil {
		utils.Logger.Warnf("validation failed: %v", err)
		return err
	}
	result := r.DB.WithContext(ctx).Save(movie)
	if result.Error != nil {
		utils.Logger.Errorf("%s ID %d: %v", constants.ErrMsgUpdateMovie, movie.ID, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		utils.Logger.Warnf("%s: ID %d", constants.MsgMovieNotFound, movie.ID)
		return ErrMovieNotFound
	}
	utils.Logger.Infof("%s: ID %d", constants.MsgMovieUpdatedSuccessfully, movie.ID)
	return nil
}

func (r *MovieRepo) GetByID(ctx context.Context, id uint) (*domain.Movie, error) {
	var movie domain.Movie
	err := r.DB.WithContext(ctx).First(&movie, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warnf("%s: ID %d", constants.MsgMovieNotFound, id)
			return nil, ErrMovieNotFound
		}
		utils.Logger.Errorf("%s %d: %v", constants.ErrMsgGetMovieByID, id, err)
		return nil, err
	}
	utils.Logger.Infof("%s: ID %d", constants.MsgMovieRetrievedSuccessfully, id)
	return &movie, nil
}

func (r *MovieRepo) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&domain.Movie{}, id)
	if result.Error != nil {
		utils.Logger.Errorf("%s ID %d: %v", constants.ErrMsgDeleteMovie, id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		utils.Logger.Warnf("%s: ID %d", constants.MsgMovieNotFound, id)
		return ErrMovieNotFound
	}
	utils.Logger.Infof("%s: ID %d", constants.MsgMovieDeletedSuccessfully, id)
	return nil
}

func (r *MovieRepo) UpdateFields(ctx context.Context, id uint, updates map[string]interface{}) error {
	result := r.DB.WithContext(ctx).Model(&domain.Movie{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		utils.Logger.Errorf("%s ID %d: %v", constants.ErrMsgUpdateMovieFields, id, result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		utils.Logger.Warnf("%s: ID %d", constants.MsgMovieNotFound, id)
		return ErrMovieNotFound
	}
	utils.Logger.Infof("%s ID %d", constants.MsgMovieUpdatedSuccessfully, id)
	return nil
}
