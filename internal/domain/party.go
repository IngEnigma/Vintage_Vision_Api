package domain

import "time"

type Party struct {
	ID        uint   `gorm:"primaryKey"`
	HostID    uint   `gorm:"not null"`
	MovieID   uint   `gorm:"not null"`
	PartyCode string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	ExpiresAt time.Time

	Movie   Movie         `gorm:"foreignKey:MovieID"`
	Host    User          `gorm:"foreignKey:HostID"`
	Members []PartyMember `gorm:"foreignKey:PartyID"`
}
