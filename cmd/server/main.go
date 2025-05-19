package main

import (
	"vintage-vision-api/infrastructure/config"
	"vintage-vision-api/infrastructure/db"
	"vintage-vision-api/infrastructure/router"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/handler"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"
	"vintage-vision-api/repository"

	"github.com/joho/godotenv"

	_ "vintage-vision-api/docs"
)

// Package main contiene la documentación principal de la API.
// @title Vintage Vision API
// @version 1.0
// @description API para la plataforma de streaming de películas clásicas.
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	if err := godotenv.Load(); err != nil {
		utils.Logger.Warn(constants.ErrMsgMissingEnvVar)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Logger.Fatalf("Error al cargar configuración: %v", err)
	}

	if err := config.InitCloudinary(); err != nil {
		utils.Logger.Fatalf("%s: %v", constants.ErrMsgInitCloudinary, err)
	}
	utils.Logger.Info(constants.MsgCloudinarySuccess)

	dbConn, err := db.Connect(cfg)
	if err != nil {
		utils.Logger.Fatalf("%s: %v", constants.ErrMsgDBConnection, err)
	}
	utils.Logger.Info(constants.MsgDBConnectionSuccess)

	userRepo := repository.NewUserRepo(dbConn.DB)
	jwtGen := &utils.JWTGenerator{}
	userUC := usecase.NewUserUsecase(userRepo, jwtGen)
	authHandler := handler.NewAuthHandler(userUC)

	profileRepo := repository.NewProfileRepo(dbConn.DB)
	profileUC := usecase.NewProfileUsecase(profileRepo)
	profileHandler := handler.NewProfileHandler(profileUC)

	movieRepo := repository.NewMovieRepo(dbConn.DB)
	movieUC := usecase.NewMovieUsecase(movieRepo)
	movieHandler := handler.NewMovieHandler(movieUC)

	r := router.SetupRouter(authHandler, profileHandler, movieHandler)

	utils.Logger.Info(constants.MsgServerStartedSuccessfully)
	if err := r.Run(":8080"); err != nil {
		utils.Logger.Fatalf("%s: %v", constants.ErrMsgServerStart, err)
	}
}
