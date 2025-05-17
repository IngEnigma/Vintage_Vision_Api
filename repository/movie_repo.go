package repository

import (
	"context"
	"errors"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"
	"vintage-vision-api/repository/entity"
	"vintage-vision-api/repository/mapper"

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
	var entities []entity.Movie
	err := r.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&entities).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetMovies, err)
		return nil, err
	}

	domainMovies := make([]domain.Movie, len(entities))
	for i, m := range entities {
		domainMovies[i] = *mapper.ToDomainMovie(&m)
	}

	utils.Logger.Infof("%s", constants.MsgMoviesRetrievedSuccessfully)
	return domainMovies, nil
}

func (r *MovieRepo) Create(ctx context.Context, movie *domain.Movie) error {
	entityMovie := mapper.FromDomainMovie(movie)
	err := r.DB.WithContext(ctx).Create(entityMovie).Error
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCreateMovie, err)
		return err
	}

	movie.ID = entityMovie.ID
	utils.Logger.Infof("%s: ID %d", constants.MsgMovieCreatedSuccessfully, movie.ID)
	return nil
}

func (r *MovieRepo) Update(ctx context.Context, movie *domain.Movie) error {
	entityMovie := mapper.FromDomainMovie(movie)

	result := r.DB.WithContext(ctx).
		Model(&entityMovie).
		Where("id = ?", movie.ID).
		Omit("WatchHistories").
		Updates(entityMovie)

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
	var entityMovie entity.Movie
	err := r.DB.WithContext(ctx).First(&entityMovie, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warnf("%s: ID %d", constants.MsgMovieNotFound, id)
			return nil, err
		}
		utils.Logger.Errorf("%s %d: %v", constants.ErrMsgGetMovieByID, id, err)
		return nil, err
	}
	domainMovie := mapper.ToDomainMovie(&entityMovie)
	utils.Logger.Infof("%s: ID %d", constants.MsgMovieRetrievedSuccessfully, id)
	return domainMovie, nil
}

func (r *MovieRepo) Delete(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&entity.Movie{}, id)
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
	result := r.DB.WithContext(ctx).Model(&entity.Movie{}).Where("id = ?", id).Updates(updates)
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
