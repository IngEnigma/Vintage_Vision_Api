package usecase

import (
	"context"
	"errors"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/utils"
)

type MovieUsecase struct {
	Repo domain.MovieRepository
}

func NewMovieUsecase(r domain.MovieRepository) *MovieUsecase {
	return &MovieUsecase{Repo: r}
}

func (u *MovieUsecase) GetAll(ctx context.Context, limit, offset int) ([]domain.Movie, error) {
	movies, err := u.Repo.GetAll(ctx, limit, offset)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetMovies, err)
		return nil, err
	}
	utils.Logger.Infof("%s, total: %d", constants.MsgMoviesRetrievedSuccessfully, len(movies))
	return movies, nil
}

func (u *MovieUsecase) Create(ctx context.Context, req request.CreateMovieRequest) (*domain.Movie, error) {
	movie := domain.Movie{
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		ImageURL:    req.ImageURL,
		StreamURL:   req.StreamURL,
		Genre:       req.Genre,
		Duration:    req.Duration,
	}

	if err := u.Repo.Create(ctx, &movie); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgCreateMovie, err)
		return nil, err
	}

	utils.Logger.Infof("%s: '%s'", constants.MsgMovieCreatedSuccessfully, movie.Title)
	return &movie, nil
}

func (u *MovieUsecase) Update(ctx context.Context, id uint, req request.UpdateMovieRequest) (*domain.Movie, error) {
	updates := map[string]interface{}{}

	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.StreamURL != nil {
		updates["stream_url"] = *req.StreamURL
	}
	if req.Genre != nil {
		updates["genre"] = *req.Genre
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}

	if len(updates) == 0 {
		utils.Logger.Warnf("No se recibieron campos para actualizar la película ID %d", id)
		return nil, errors.New("no se proporcionaron datos para actualizar")
	}

	if err := u.Repo.UpdateFields(ctx, id, updates); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgUpdateMovie, err)
		return nil, err
	}

	updatedMovie, err := u.Repo.GetByID(ctx, id)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetUpdatedMovie, err)
		return nil, err
	}

	utils.Logger.Infof("%s con ID %d", constants.MsgMovieUpdatedSuccessfully, id)
	return updatedMovie, nil
}

func (u *MovieUsecase) Delete(ctx context.Context, id uint) error {
	if err := u.Repo.Delete(ctx, id); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgDeleteMovie, err)
		return err
	}
	utils.Logger.Infof("%s con ID %d", constants.MsgMovieDeletedSuccessfully, id)
	return nil
}
