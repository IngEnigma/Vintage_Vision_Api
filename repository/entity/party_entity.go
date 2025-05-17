package entity

import (
	"time"

	"gorm.io/gorm"
)

type Party struct {
	gorm.Model
	HostID    uint   `gorm:"not null"`
	MovieID   uint   `gorm:"not null"`
	PartyCode string `gorm:"uniqueIndex"`
	ExpiresAt time.Time

	Movie   Movie         `gorm:"foreignKey:MovieID"`
	Host    User          `gorm:"foreignKey:HostID"`
	Members []PartyMember `gorm:"foreignKey:PartyID"`
}
