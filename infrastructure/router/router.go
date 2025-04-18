package router

import (
	"vintage-vision-api/internal/handler"
	"vintage-vision-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, movieHandler *handler.MovieHandler) *gin.Engine {
	r := gin.Default()

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
	private.POST("/profiles", profileHandler.Create)
	private.GET("/profiles", profileHandler.GetAll)
	private.DELETE("/profiles/:id", profileHandler.Delete)
	private.PATCH("/profiles/:id", profileHandler.Update)
	private.GET("/movies", movieHandler.GetAll)
}

func adminRoutes(admin *gin.RouterGroup, movieHandler *handler.MovieHandler) {
	admin.POST("/movies", movieHandler.Create)
	admin.PATCH("/movies/:id", movieHandler.Update)
	admin.DELETE("/movies/:id", movieHandler.Delete)
}
