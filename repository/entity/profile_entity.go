package entity

import (
	"errors"
	"vintage-vision-api/internal/constants"

	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	Name      string `gorm:"not null"`
	AvatarURL string
	UserID    uint `gorm:"not null;index"`

	WatchHistory []WatchHistory `gorm:"foreignKey:ProfileID"`
}

func (p *Profile) Validate() error {
	if p.Name == "" {
		return errors.New(constants.ErrMshEmptyProfileName)
	}
	if p.UserID == 0 {
		return errors.New(constants.ErrMsgEmptyProfileId)
	}
	return nil
}
