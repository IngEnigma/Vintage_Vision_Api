package domain

import "time"

type WatchHistory struct {
	ID        uint
	ProfileID uint
	MovieID   uint
	WatchedAt time.Time
	Progress  float32 // % of the movie watched 0.0 to 1.0
}
