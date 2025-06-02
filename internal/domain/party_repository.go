package domain

import "context"

type PartyRepository interface {
	Create(ctx context.Context, party *Party) error
	FindByCode(ctx context.Context, code string) (*Party, error)
	AddMember(ctx context.Context, partyID, profileID uint) error
	RemoveMember(ctx context.Context, partyID, profileID uint) error
	IsMember(ctx context.Context, partyID, profileID uint) (bool, error)
}
