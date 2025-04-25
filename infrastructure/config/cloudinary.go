package config

import (
	"os"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/utils"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/joho/godotenv"
)

var CLD *cloudinary.Cloudinary

func InitCloudinary() error {
	if err := godotenv.Load(); err != nil {
		utils.Logger.Warn(constants.ErrMsgMissingEnvVar)
	}

	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgInitCloudinary, err)
		return err
	}

	CLD = cld
	return nil
}
