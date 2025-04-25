package config

import (
	"os"

	"fmt"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/utils"
)

type Config struct {
	DatabaseURL string
}

func LoadConfig() (*Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		err := fmt.Errorf("%s: DATABASE_URL", constants.ErrMsgMissingEnvVar)
		utils.Logger.Error(err)
		return nil, err
	}

	utils.Logger.Info(constants.MsgEnvCargedSuccessfully)

	return &Config{
		DatabaseURL: databaseURL,
	}, nil
}
