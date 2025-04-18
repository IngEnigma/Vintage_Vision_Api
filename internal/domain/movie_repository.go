package domain

type MovieRepository interface {
	GetAll() ([]Movie, error)
	Create(movie *Movie) error
	GetByID(id uint) (*Movie, error)
	Update(movie *Movie) error
	Delete(id uint) error
	UpdateFields(id uint, updates map[string]interface{}) error
}
