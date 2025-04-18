package domain

import "time"

type WatchHistory struct {
	ID        uint `gorm:"primaryKey"`
	ProfileID uint `gorm:"not null"`
	MovieID   uint `gorm:"not null"`
	WatchedAt time.Time
	Progress  float32 // % de progreso visto (0.0 - 1.0)
}
