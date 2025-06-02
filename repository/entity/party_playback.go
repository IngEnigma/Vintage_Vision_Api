package entity

import (
	"time"

	"gorm.io/gorm"
)

type PartyPlayback struct {
	gorm.Model
	PartyID   uint `gorm:"uniqueIndex"`
	IsPlaying bool
	CurrentAt float64
	UpdatedAt time.Time
}
