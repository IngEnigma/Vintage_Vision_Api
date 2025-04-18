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
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
		//api.GET("/movies", movieHandler.GetAll)

		private := api.Group("/")
		private.Use(middleware.AuthMiddleware())
		{
			private.GET("/me", func(c *gin.Context) {
				userID := c.GetUint("user_id")
				c.JSON(200, gin.H{"message": "Estás autenticado", "user_id": userID})
			})
			private.POST("/profiles", profileHandler.Create)
			private.GET("/profiles", profileHandler.GetAll)
			private.DELETE("/profiles/:id", profileHandler.Delete)
			private.PATCH("/profiles/:id", profileHandler.Update)
			private.GET("/movies", movieHandler.GetAll)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
		{
			admin.POST("/movies", movieHandler.Create)
			admin.PUT("/movies/:id", movieHandler.Update)
			admin.DELETE("/movies/:id", movieHandler.Delete)
		}

	}

	return r
}
