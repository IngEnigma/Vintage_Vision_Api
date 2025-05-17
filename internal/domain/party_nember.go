package domain

import "time"

type PartyMember struct {
	ID        uint
	PartyID   uint
	ProfileID uint
	JoinedAt  time.Time
}
