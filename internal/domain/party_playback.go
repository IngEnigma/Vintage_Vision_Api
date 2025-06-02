package domain

import "time"

type PartyPlayback struct {
	ID        uint
	PartyID   uint
	IsPlaying bool
	CurrentAt float64
	UpdatedAt time.Time
}
