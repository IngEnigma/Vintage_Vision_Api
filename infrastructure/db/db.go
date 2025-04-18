package db

import (
	"fmt"

	"vintage-vision-api/infrastructure/config"
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
		utils.Logger.Errorf("Error al conectar con la base de datos: %v", err)
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
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
		utils.Logger.Errorf("Error al hacer migraciones: %v", err)
		return nil, fmt.Errorf("error al hacer migraciones: %w", err)
	}

	utils.Logger.Info("Conectado y migraciones ejecutadas correctamente")
	return &DBConnection{db}, nil
}

func Connect(cfg *config.Config) (*DBConnection, error) {
	return NewDBConnection(cfg.DatabaseURL)
}
