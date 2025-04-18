package handler

import (
	"errors"
	"net/http"
	"strconv"

	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	Usecase *usecase.ProfileUsecase
}

func NewProfileHandler(u *usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{Usecase: u}
}

func (h *ProfileHandler) Create(c *gin.Context) {
	var req request.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userID := c.GetUint("user_id")
	if err := h.Usecase.Create(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el perfil"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Perfil creado"})
}

func (h *ProfileHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	profiles, err := h.Usecase.GetAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener perfiles"})
		return
	}

	var res []response.ProfileResponse
	for _, p := range profiles {
		res = append(res, response.ProfileResponse{
			ID:        p.ID,
			Name:      p.Name,
			AvatarUrl: p.AvatarURL,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (h *ProfileHandler) Delete(c *gin.Context) {
	profileIDStr := c.Param("id")
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userID := c.GetUint("user_id")
	err = h.Usecase.Delete(uint(profileID), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil no encontrado o no permitido"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el perfil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil eliminado"})
}

func (h *ProfileHandler) Update(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userID := c.GetUint("user_id")
	err = h.Usecase.Update(uint(profileID), userID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil no encontrado o no permitido"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el perfil"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado"})
}
