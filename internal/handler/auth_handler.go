package handler

import (
	"net/http"

	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Usecase *usecase.UserUsecase
}

func NewAuthHandler(u *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{Usecase: u}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warnf("Registro fallido: datos inválidos - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	if err := h.Usecase.Register(req); err != nil {
		utils.Logger.Errorf("Error al registrar usuario %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo registrar el usuario", "details": err.Error()})
		return
	}

	utils.Logger.Infof("Usuario registrado correctamente: %s", req.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado correctamente"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warnf("Inicio de sesión fallido: datos inválidos - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	token, err := h.Usecase.Login(req)
	if err != nil {
		utils.Logger.Warnf("Inicio de sesión fallido para %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	utils.Logger.Infof("Usuario autenticado: %s", req.Email)
	c.JSON(http.StatusOK, response.AuthResponse{Token: token})
}
