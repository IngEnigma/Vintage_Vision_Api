package domain

import "context"

type ProfileRepository interface {
	Create(ctx context.Context, profile *Profile) error
	FindByUser(ctx context.Context, userID uint) ([]Profile, error)
	DeleteByID(ctx context.Context, profileID uint, userID uint) error
	Update(ctx context.Context, profile *Profile) error
	FindByIDAndUser(ctx context.Context, profileID, userID uint) (*Profile, error)
}
