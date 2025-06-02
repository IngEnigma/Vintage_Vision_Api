package entity

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	Name      string `gorm:"not null"`
	AvatarURL string
	UserID    uint `gorm:"not null;index"`

	WatchHistory []WatchHistory `gorm:"foreignKey:ProfileID"`
}
