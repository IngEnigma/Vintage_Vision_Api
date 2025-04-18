package domain

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
}

type ProfileRepository interface {
	Create(profile *Profile) error
	FindByUser(userID uint) ([]Profile, error)
	DeleteByID(profileID uint, userID uint) error
	Update(profile *Profile) error
	FindByIDAndUser(profileID, userID uint) (*Profile, error)
}

type MovieRepository interface {
	GetAll() ([]Movie, error)
	Create(movie *Movie) error
	GetByID(id uint) (*Movie, error)
	Update(movie *Movie) error
	Delete(id uint) error
}
