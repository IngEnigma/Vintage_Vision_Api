package domain

import "time"

type Movie struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null;size:200"`
	Description string
	Year        int
	Genre       string
	ImageURL    string
	StreamURL   string `gorm:"not null;size:500"`
	Duration    int
	CreatedAt   time.Time
	UpdatedAt   time.Time

	WatchHistories []WatchHistory `gorm:"foreignKey:MovieID"`
}
