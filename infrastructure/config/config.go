package config

import (
	"os"

	"vintage-vision-api/internal/utils"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		utils.Logger.Fatal("DATABASE_URL no está definido en el entorno")
	}

	utils.Logger.Info("Configuración cargada correctamente desde variables de entorno")

	return &Config{
		DatabaseURL: databaseURL,
	}
}
