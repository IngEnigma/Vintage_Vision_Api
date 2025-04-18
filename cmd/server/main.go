package main

import (
	"log"

	"github.com/joho/godotenv"

	"vintage-vision-api/infrastructure/db"
	"vintage-vision-api/infrastructure/router"
	"vintage-vision-api/internal/handler"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	db.Connect()

	userRepo := repository.NewUserRepo(db.DB)
	userUC := usecase.NewUserUsecase(userRepo)
	authHandler := handler.NewAuthHandler(userUC)

	profileRepo := repository.NewProfileRepo(db.DB)
	profileUC := usecase.NewProfileUsecase(profileRepo)
	profileHandler := handler.NewProfileHandler(profileUC)

	movieRepo := repository.NewMovieRepo(db.DB)
	movieUC := usecase.NewMovieUsecase(movieRepo)
	movieHandler := handler.NewMovieHandler(movieUC)

	r := router.SetupRouter(authHandler, profileHandler, movieHandler)
	r.Run(":8080")

}
