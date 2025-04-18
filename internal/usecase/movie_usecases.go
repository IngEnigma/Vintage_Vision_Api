package usecase

import (
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

func (u *MovieUsecase) GetAll() ([]domain.Movie, error) {
	movies, err := u.Repo.GetAll()
	if err != nil {
		utils.Logger.Errorf("Error al obtener las películas: %v", err)
		return nil, err
	}
	utils.Logger.Infof("Películas obtenidas correctamente, total: %d", len(movies))
	return movies, nil
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

	err := u.Repo.Create(&movie)
	if err != nil {
		utils.Logger.Errorf("Error al crear la película '%s': %v", req.Title, err)
		return err
	}
	utils.Logger.Infof("Película '%s' creada con éxito", req.Title)
	return nil
}

func (u *MovieUsecase) Update(id uint, req request.UpdateMovieRequest) error {
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
		utils.Logger.Warnf("No se recibieron campos para actualizar en la película con ID %d", id)
		return nil
	}

	err := u.Repo.UpdateFields(id, updates)
	if err != nil {
		utils.Logger.Errorf("Error al actualizar la película con ID %d: %v", id, err)
		return err
	}

	utils.Logger.Infof("Película con ID %d actualizada con éxito", id)
	return nil
}

func (u *MovieUsecase) Delete(id uint) error {
	err := u.Repo.Delete(id)
	if err != nil {
		utils.Logger.Errorf("Error al eliminar la película con ID %d: %v", id, err)
		return err
	}
	utils.Logger.Infof("Película con ID %d eliminada correctamente", id)
	return nil
}
