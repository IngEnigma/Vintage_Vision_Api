package entity

import (
	"time"

	"gorm.io/gorm"
)

type PartyMember struct {
	gorm.Model
	PartyID   uint `gorm:"not null"`
	ProfileID uint `gorm:"not null"`
	JoinedAt  time.Time
}
