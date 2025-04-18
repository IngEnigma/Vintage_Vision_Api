package domain

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null;size:255"`
	Password  string `gorm:"not null"`
	IsAdmin   bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Profiles []Profile `gorm:"foreignKey:UserID"`
}

func (u *User) IsAdminUser() bool {
	return u.IsAdmin
}
