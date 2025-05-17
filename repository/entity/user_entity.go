package entity

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null;size:255"`
	Password string `gorm:"not null"`
	IsAdmin  bool   `gorm:"default:false"`

	Profiles []Profile `gorm:"foreignKey:UserID"`
}

func (u *User) IsAdminUser() bool {
	return u.IsAdmin
}
