package domain

import "time"

type WatchHistory struct {
	ID        uint `gorm:"primaryKey"`
	ProfileID uint `gorm:"not null;uniqueIndex:idx_profile_movie"`
	MovieID   uint `gorm:"not null;uniqueIndex:idx_profile_movie"`
	WatchedAt time.Time
	Progress  float32 // % of the movie watched 0.0 to 1.0
}
