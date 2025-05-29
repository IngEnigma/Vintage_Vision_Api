package router

import (
	"time"
	"vintage-vision-api/internal/handler"
	"vintage-vision-api/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, movieHandler *handler.MovieHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		publicRoutes(api, authHandler)
		private := api.Group("/")
		private.Use(middleware.AuthMiddleware())
		privateRoutes(private, profileHandler, movieHandler)

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
		adminRoutes(admin, movieHandler)
	}

	return r
}

func publicRoutes(api *gin.RouterGroup, authHandler *handler.AuthHandler) {
	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)
}

func privateRoutes(private *gin.RouterGroup, profileHandler *handler.ProfileHandler, movieHandler *handler.MovieHandler) {
	private.GET("/movies/genre", movieHandler.GetMoviesByGenre)
	private.DELETE("/profiles/:id", profileHandler.Delete)
	private.PATCH("/profiles/:id", profileHandler.Update)
	private.GET("/movies/:id", movieHandler.GetByID)
	private.POST("/profiles", profileHandler.Create)
	private.GET("/profiles", profileHandler.GetAll)
	private.GET("/movies", movieHandler.GetAll)
	private.GET("/movies/detail/:id", movieHandler.GetMovieDetail)
	private.GET("/movies/preview/:id", movieHandler.GetPreviewByGenre)
	private.GET("/movies/player/:id", movieHandler.GetMoviePlayer)
}

func adminRoutes(admin *gin.RouterGroup, movieHandler *handler.MovieHandler) {
	admin.DELETE("/movies/:id", movieHandler.Delete)
	admin.PATCH("/movies/:id", movieHandler.Update)
	admin.POST("/movies", movieHandler.Create)
}
