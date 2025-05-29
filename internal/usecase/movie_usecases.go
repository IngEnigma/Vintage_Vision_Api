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

func (uc *MovieUsecase) GetMoviesByGenre(ctx context.Context, genre string, page int, limit int) ([]domain.Movie, error) {
	if genre == "" {
		utils.Logger.Warnf("%s: %s", constants.ErrMsgInvalidGenre, genre)
		return nil, errors.New(constants.ErrMsgInvalidGenre)
	}

	if !utils.IsValidGenre(genre) {
		utils.Logger.Warnf("%s: %s", constants.ErrMsgInvalidGenre, genre)
		return nil, errors.New(constants.ErrMsgInvalidGenre)
	}

	page, limit = utils.ValidatePaginationParams(page, limit)
	return uc.Repo.GetMoviesByGenre(ctx, genre, page, limit)
}

func (u *MovieUsecase) Create(ctx context.Context, req request.CreateMovieRequest) (*domain.Movie, error) {
	validationErrors := utils.ValidateCreateMovieInput(req)
	if len(validationErrors) > 0 {
		utils.Logger.Warnf("Errores de validación al crear película: %+v", validationErrors)
		return nil, errors.New(constants.ErrMsgInvalidInput)
	}

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
	validationErrors := utils.ValidateUpdateMovieInput(req)
	if len(validationErrors) > 0 {
		utils.Logger.Warnf("Errores de validación al actualizar película %d: %+v", id, validationErrors)
		return nil, errors.New(constants.ErrMsgInvalidInput)
	}

	updates := utils.BuildMovieUpdateMap(req)

	if len(updates) == 0 {
		utils.Logger.Warnf("%s: %d", constants.ErrMsgUpdateMovieFields, id)
		return nil, errors.New(constants.ErrMsgUpdateMovieFields)
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
	movie, err := u.Repo.GetByID(ctx, id)
	if err != nil || movie == nil {
		utils.Logger.Warnf("%s: ID %d", constants.ErrMsgMovieNotFound, id)
		return errors.New(constants.ErrMsgMovieNotFound)
	}

	if err := u.Repo.Delete(ctx, id); err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgDeleteMovie, err)
		return err
	}

	utils.Logger.Infof("%s con ID %d", constants.MsgMovieDeletedSuccessfully, id)
	return nil
}

func (u *MovieUsecase) GetByID(ctx context.Context, id uint) (*domain.Movie, error) {
	movie, err := u.Repo.GetByID(ctx, id)
	if err != nil || movie == nil {
		utils.Logger.Warnf("%s: ID %d", constants.ErrMsgMovieNotFound, id)
		return nil, errors.New(constants.ErrMsgMovieNotFound)
	}

	utils.Logger.Infof("%s con ID %d", constants.MsgMovieRetrievedSuccessfully, id)
	return movie, nil
}
