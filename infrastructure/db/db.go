package db

import (
	"fmt"

	"vintage-vision-api/infrastructure/config"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnection struct {
	*gorm.DB
}

func NewDBConnection(dsn string) (*DBConnection, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgDBConnection, err)
		return nil, fmt.Errorf("%s: %w", constants.ErrMsgDBConnection, err)
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.Movie{},
		&domain.WatchHistory{},
		&domain.Party{},
		&domain.PartyMember{},
	)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgDBMigration, err)
		return nil, fmt.Errorf("%s: %w", constants.ErrMsgDBMigration, err)
	}

	utils.Logger.Info(constants.MsgDBMigrationSuccess)
	return &DBConnection{db}, nil
}

func Connect(cfg *config.Config) (*DBConnection, error) {
	return NewDBConnection(cfg.DatabaseURL)
}
