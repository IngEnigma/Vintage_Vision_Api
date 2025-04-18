package repository

import (
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"

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
	if err != nil {
		utils.Logger.Errorf("Error al obtener todas las películas: %v", err)
		return nil, err
	}
	utils.Logger.Infof("Películas obtenidas correctamente, total: %d", len(movies))
	return movies, nil
}

func (r *MovieRepo) Create(movie *domain.Movie) error {
	err := r.DB.Create(movie).Error
	if err != nil {
		utils.Logger.Errorf("Error al crear la película '%s': %v", movie.Title, err)
		return err
	}
	utils.Logger.Infof("Película '%s' creada correctamente", movie.Title)
	return nil
}

func (r *MovieRepo) GetByID(id uint) (*domain.Movie, error) {
	var movie domain.Movie
	err := r.DB.First(&movie, id).Error
	if err != nil {
		utils.Logger.Warnf("No se encontró la película con ID %d: %v", id, err)
		return nil, err
	}
	utils.Logger.Infof("Película con ID %d obtenida correctamente", id)
	return &movie, nil
}

func (r *MovieRepo) Update(movie *domain.Movie) error {
	err := r.DB.Save(movie).Error
	if err != nil {
		utils.Logger.Errorf("Error al actualizar la película '%s': %v", movie.Title, err)
		return err
	}
	utils.Logger.Infof("Película '%s' actualizada correctamente", movie.Title)
	return nil
}

func (r *MovieRepo) Delete(id uint) error {
	err := r.DB.Delete(&domain.Movie{}, id).Error
	if err != nil {
		utils.Logger.Errorf("Error al eliminar la película con ID %d: %v", id, err)
		return err
	}
	utils.Logger.Infof("Película con ID %d eliminada correctamente", id)
	return nil
}

func (r *MovieRepo) UpdateFields(id uint, updates map[string]interface{}) error {
	return r.DB.Model(&domain.Movie{}).Where("id = ?", id).Updates(updates).Error
}
