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

// @description Configura las rutas públicas, privadas y de administración
func SetupRouter(authHandler *handler.AuthHandler, profileHandler *handler.ProfileHandler, movieHandler *handler.MovieHandler) *gin.Engine {
	r := gin.Default()

	// Configuración CORS más segura para producción (ajusta los orígenes según necesites)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // En producción cambia a tus dominios específicos
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes grouping
	api := r.Group("/api")
	{
		publicRoutes(api, authHandler) // Public routes (no auth required)

		private := api.Group("/") // Private routes (auth required)
		private.Use(middleware.AuthMiddleware())
		privateRoutes(private, profileHandler, movieHandler)

		admin := api.Group("/admin") // Admin routes (auth + admin role required)
		admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
		adminRoutes(admin, movieHandler)
	}

	return r
}

// publicRoutes configura las rutas accesibles sin autenticación
// @Summary Rutas públicas
// @Description Endpoints que no requieren autenticación
func publicRoutes(api *gin.RouterGroup, authHandler *handler.AuthHandler) {
	// Registro de usuario
	// swagger:route POST /api/register auth registerUser
	//
	// Crea una nueva cuenta de usuario
	//
	// Responses:
	//  201: SuccessResponse
	//  400: ErrorResponse
	//  409: ErrorResponse
	//  500: ErrorResponse
	api.POST("/register", authHandler.Register)

	// Inicio de sesión
	// swagger:route POST /api/login auth loginUser
	//
	// Autentica un usuario y devuelve un token JWT
	//
	// Responses:
	//  200: AuthResponse
	//  400: ErrorResponse
	//  401: ErrorResponse
	//  500: ErrorResponse
	api.POST("/login", authHandler.Login)
}

// privateRoutes configura las rutas que requieren autenticación
// @Summary Rutas privadas
// @Description Endpoints que requieren autenticación JWT válida
func privateRoutes(private *gin.RouterGroup, profileHandler *handler.ProfileHandler, movieHandler *handler.MovieHandler) {
	// Perfiles
	// swagger:route POST /api/profiles profiles createProfile
	//
	// Crea un nuevo perfil para el usuario autenticado
	//
	// Security:
	// - BearerAuth: []
	//
	// Responses:
	//  201: ProfileResponse
	//  400: ErrorResponse
	//  401: ErrorResponse
	//  500: ErrorResponse
	private.POST("/profiles", profileHandler.Create)

	// swagger:route GET /api/profiles profiles listProfiles
	//
	// Obtiene todos los perfiles del usuario autenticado
	//
	// Security:
	// - BearerAuth: []
	//
	// Responses:
	//  200: []ProfileResponse
	//  401: ErrorResponse
	//  500: ErrorResponse
	private.GET("/profiles", profileHandler.GetAll)

	// swagger:route DELETE /api/profiles/{id} profiles deleteProfile
	//
	// Elimina un perfil específico
	//
	// Security:
	// - BearerAuth: []
	//
	// Parameters:
	// + name: id
	//   in: path
	//   description: ID del perfil
	//   required: true
	//   type: string
	//
	// Responses:
	//  204: empty
	//  401: ErrorResponse
	//  404: ErrorResponse
	//  500: ErrorResponse
	private.DELETE("/profiles/:id", profileHandler.Delete)

	// swagger:route PATCH /api/profiles/{id} profiles updateProfile
	//
	// Actualiza un perfil específico
	//
	// Security:
	// - BearerAuth: []
	//
	// Parameters:
	// + name: id
	//   in: path
	//   description: ID del perfil
	//   required: true
	//   type: string
	//
	// Responses:
	//  200: ProfileResponse
	//  400: ErrorResponse
	//  401: ErrorResponse
	//  404: ErrorResponse
	//  500: ErrorResponse
	private.PATCH("/profiles/:id", profileHandler.Update)

	// Películas
	// swagger:route GET /api/movies movies listMovies
	//
	// Obtiene el catálogo de películas disponibles
	//
	// Security:
	// - BearerAuth: []
	//
	// Responses:
	//  200: []MovieResponse
	//  401: ErrorResponse
	//  500: ErrorResponse
	private.GET("/movies", movieHandler.GetAll)
}

// adminRoutes configura las rutas exclusivas para administradores
// @Summary Rutas de administrador
// @Description Endpoints que requieren autenticación y rol de administrador
func adminRoutes(admin *gin.RouterGroup, movieHandler *handler.MovieHandler) {
	// swagger:route POST /api/admin/movies movies createMovie
	//
	// Crea una nueva película (solo administradores)
	//
	// Security:
	// - BearerAuth: []
	//
	// Responses:
	//  201: MovieResponse
	//  400: ErrorResponse
	//  401: ErrorResponse
	//  403: ErrorResponse
	//  500: ErrorResponse
	admin.POST("/movies", movieHandler.Create)

	// swagger:route PATCH /api/admin/movies/{id} movies updateMovie
	//
	// Actualiza una película existente (solo administradores)
	//
	// Security:
	// - BearerAuth: []
	//
	// Parameters:
	// + name: id
	//   in: path
	//   description: ID de la película
	//   required: true
	//   type: string
	//
	// Responses:
	//  200: MovieResponse
	//  400: ErrorResponse
	//  401: ErrorResponse
	//  403: ErrorResponse
	//  404: ErrorResponse
	//  500: ErrorResponse
	admin.PATCH("/movies/:id", movieHandler.Update)

	// swagger:route DELETE /api/admin/movies/{id} movies deleteMovie
	//
	// Elimina una película (solo administradores)
	//
	// Security:
	// - BearerAuth: []
	//
	// Parameters:
	// + name: id
	//   in: path
	//   description: ID de la película
	//   required: true
	//   type: string
	//
	// Responses:
	//  204: empty
	//  401: ErrorResponse
	//  403: ErrorResponse
	//  404: ErrorResponse
	//  500: ErrorResponse
	admin.DELETE("/movies/:id", movieHandler.Delete)
}
