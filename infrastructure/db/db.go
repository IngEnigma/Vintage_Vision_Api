package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"vintage-vision-api/internal/domain"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("🔴 DATABASE_URL no está definido en el entorno")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("🔴 Error al conectar con la base de datos: %v", err)
	}

	log.Println("🟢 Conectado a la base de datos")

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Profile{},
		&domain.Movie{},
		&domain.WatchHistory{},
		&domain.Party{},
		&domain.PartyMember{},
	)

	if err != nil {
		log.Fatalf("🔴 Error al hacer migraciones: %v", err)
	}

	DB = db
}
