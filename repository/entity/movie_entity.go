package entity

import (
	"errors"
	"strings"

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

func (m *Movie) Validate() error {
	if strings.TrimSpace(m.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(m.StreamURL) == "" {
		return errors.New("stream URL is required")
	}
	return nil
}
