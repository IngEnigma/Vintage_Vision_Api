package domain

import "time"

type Profile struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"`
	AvatarURL string
	UserID    uint `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	WatchHistory []WatchHistory
}
