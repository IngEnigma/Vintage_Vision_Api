package domain

import "time"

type Movie struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Year        int
	Genre       string
	PosterURL   string
	VideoURL    string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
