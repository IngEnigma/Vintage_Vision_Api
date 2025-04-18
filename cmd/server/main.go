package main

import (
	"vintage-vision-api/infrastructure/config"
	"vintage-vision-api/infrastructure/db"
	"vintage-vision-api/infrastructure/router"
	"vintage-vision-api/internal/handler"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"
	"vintage-vision-api/repository"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		utils.Logger.Warn("No se encontró el archivo .env, usando variables de entorno del sistema")
	}

	cfg := config.LoadConfig()

	dbConn, err := db.Connect(cfg)
	if err != nil {
		utils.Logger.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	utils.Logger.Info("Conexión a la base de datos exitosa")

	userRepo := repository.NewUserRepo(dbConn.DB)
	userUC := usecase.NewUserUsecase(userRepo)
	authHandler := handler.NewAuthHandler(userUC)

	profileRepo := repository.NewProfileRepo(dbConn.DB)
	profileUC := usecase.NewProfileUsecase(profileRepo)
	profileHandler := handler.NewProfileHandler(profileUC)

	movieRepo := repository.NewMovieRepo(dbConn.DB)
	movieUC := usecase.NewMovieUsecase(movieRepo)
	movieHandler := handler.NewMovieHandler(movieUC)

	r := router.SetupRouter(authHandler, profileHandler, movieHandler)

	utils.Logger.Info("Servidor iniciado en http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		utils.Logger.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
