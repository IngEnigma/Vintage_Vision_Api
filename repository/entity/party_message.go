package entity

import (
	"time"

	"gorm.io/gorm"
)

type PartyMessage struct {
	gorm.Model
	PartyID  uint   `gorm:"not null"`
	SenderID uint   `gorm:"not null"`
	Message  string `gorm:"type:text"`
	SentAt   time.Time

	Sender Profile `gorm:"foreignKey:SenderID"`
}
