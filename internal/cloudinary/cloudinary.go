package cloudinary

import (
	"context"
	"errors"
	"mime/multipart"
	"time"
	"vintage-vision-api/infrastructure/config"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/utils"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	imageUploadTimeout = 60 * time.Second
	videoUploadTimeout = 5 * time.Minute
	publicFolder       = "movies/"
)

func UploadImageToCloudinary(file multipart.File, fileName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), imageUploadTimeout)
	defer cancel()

	overwrite := true
	uploadParams := uploader.UploadParams{
		ResourceType: "image",
		PublicID:     publicFolder + fileName,
		Overwrite:    &overwrite,
	}

	uploadResp, err := config.CLD.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		utils.Logger.Errorf("%s (imagen): %v", constants.ErrMsgUploadToCloudinary, err)
		return "", errors.New(constants.ErrMsgUploadToCloudinary)
	}
	return uploadResp.SecureURL, nil
}

func UploadVideoToCloudinary(file multipart.File, fileName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), videoUploadTimeout)
	defer cancel()

	overwrite := true
	uploadParams := uploader.UploadParams{
		ResourceType: "video",
		PublicID:     publicFolder + fileName,
		Overwrite:    &overwrite,
	}

	uploadResp, err := config.CLD.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		utils.Logger.Errorf("%s (video): %v", constants.ErrMsgUploadToCloudinary, err)
		return "", errors.New(constants.ErrMsgUploadToCloudinary)
	}
	return uploadResp.SecureURL, nil
}

func UploadFileToCloudinary(file *multipart.FileHeader, fileType string) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgOpenUploadedFile, err)
		return "", errors.New(constants.ErrMsgOpenUploadedFile)
	}
	defer fileContent.Close()

	switch fileType {
	case "image":
		return UploadImageToCloudinary(fileContent, file.Filename)
	case "video":
		return UploadVideoToCloudinary(fileContent, file.Filename)
	default:
		utils.Logger.Errorf("%s: %s", constants.ErrMsgUnsupportedFileType, fileType)
		return "", errors.New(constants.ErrMsgUnsupportedFileType)
	}
}
