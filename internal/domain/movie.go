package domain

import "time"

type Movie struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Year        int
	Genre       string
	ImageURL    string
	StreamURL   string `gorm:"not null"`
	Duration    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
