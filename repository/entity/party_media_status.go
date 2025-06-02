package entity

import (
	"time"

	"gorm.io/gorm"
)

type PartyMediaStatus struct {
	gorm.Model
	PartyID   uint `gorm:"not null"`
	ProfileID uint `gorm:"not null"`
	CameraOn  bool
	MicOn     bool
	UpdatedAt time.Time
}
