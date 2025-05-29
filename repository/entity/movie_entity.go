package entity

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title       string `gorm:"not null;size:200"`
	Description string
	Year        int
	Genre       string
	ImageURL    string
	StreamURL   string `gorm:"not null;size:500"`
	Duration    int

	WatchHistories []WatchHistory `gorm:"foreignKey:MovieID"`
}
