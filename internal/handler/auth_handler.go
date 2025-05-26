package handler

import (
	"net/http"

	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Usecase *usecase.UserUsecase
}

func NewAuthHandler(userUsecase *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{
		Usecase: userUsecase,
	}
}

// Register godoc
// @Summary Registrar un nuevo usuario
// @Description Crea una nueva cuenta de usuario en el sistema
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "Datos de registro"
// @Success 201 {object} response.SuccessResponse "Usuario registrado exitosamente"
// @Failure 400 {object} response.ErrorResponse "Solicitud inválida"
// @Failure 409 {object} response.ErrorResponse "Usuario ya existe"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidRequest, err)
		return
	}

	if err := h.Usecase.Register(ctx, req); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgRegisterUser, err)
		return
	}

	utils.Logger.Infof(constants.MsgUserRegisteredSuccessfully)
	c.JSON(http.StatusCreated, gin.H{"message": constants.MsgUserRegisteredSuccessfully})
}

// Login godoc
// @Summary Iniciar sesión
// @Description Autentica un usuario y devuelve un token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "Credenciales de acceso"
// @Success 200 {object} response.AuthResponse "Token JWT generado"
// @Failure 400 {object} response.ErrorResponse "Solicitud inválida"
// @Failure 401 {object} response.ErrorResponse "Credenciales inválidas"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidRequest, err)
		return
	}

	token, err := h.Usecase.Login(ctx, req)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, constants.ErrMsgInvalidCredentials, err)
		return
	}

	utils.Logger.Infof(constants.MsgUserLoggedInSuccessfully)
	c.JSON(http.StatusOK, response.AuthResponse{Token: token})
}
