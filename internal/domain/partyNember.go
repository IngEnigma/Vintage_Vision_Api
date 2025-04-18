package domain

import "time"

type PartyMember struct {
	ID        uint `gorm:"primaryKey"`
	PartyID   uint `gorm:"not null"`
	ProfileID uint `gorm:"not null"`
	JoinedAt  time.Time
}
