package domain

import "time"

type Party struct {
	ID        uint
	HostID    uint
	MovieID   uint
	PartyCode string
	CreatedAt time.Time
	ExpiresAt time.Time

	Movie   Movie
	Host    User
	Members []PartyMember
}
