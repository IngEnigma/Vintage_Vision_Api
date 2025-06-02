package domain

import "time"

type PartyMediaStatus struct {
	ID        uint
	PartyID   uint
	ProfileID uint
	CameraOn  bool
	MicOn     bool
	UpdatedAt time.Time
}
