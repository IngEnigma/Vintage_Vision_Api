package entity

import (
	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	Title       string `gorm:"not null;size:200"`
	Description string `gorm:"not null"`
	Year        int    `gorm:"not null"`
	Genre       string `gorm:"not null"`
	ImageURL    string `gorm:"not null"`
	StreamURL   string `gorm:"not null"`
	Duration    int    `gorm:"not null"`

	WatchHistories []WatchHistory `gorm:"foreignKey:MovieID"`
}
