package domain

type ProfileRepository interface {
	Create(profile *Profile) error
	FindByUser(userID uint) ([]Profile, error)
	DeleteByID(profileID uint, userID uint) error
	Update(profile *Profile) error
	FindByIDAndUser(profileID, userID uint) (*Profile, error)
}
