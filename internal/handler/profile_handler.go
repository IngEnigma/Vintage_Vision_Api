package handler

import (
	"errors"
	"net/http"
	"strconv"

	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"

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
		utils.Logger.Warnf("Error al parsear JSON en Create: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	if err := h.Usecase.Create(userID, req); err != nil {
		utils.Logger.Errorf("Error al crear perfil para usuario %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear el perfil, detalles: " + err.Error()})
		return
	}

	utils.Logger.Infof("Perfil creado para usuario %d", userID)
	c.JSON(http.StatusCreated, gin.H{"message": "Perfil creado"})
}

func (h *ProfileHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	profiles, err := h.Usecase.GetAll(userID)
	if err != nil {
		utils.Logger.Errorf("Error al obtener perfiles para usuario %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener perfiles", "details": err.Error()})
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

	utils.Logger.Infof("Perfiles obtenidos para usuario %d", userID)
	c.JSON(http.StatusOK, res)
}

func (h *ProfileHandler) Delete(c *gin.Context) {
	profileIDStr := c.Param("id")
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		utils.Logger.Warnf("ID de perfil inválido: %s", profileIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	userID := c.GetUint("user_id")
	err = h.Usecase.Delete(uint(profileID), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warnf("Perfil %d no encontrado o acceso denegado para usuario %d", profileID, userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil no encontrado o no permitido, detalles: " + err.Error()})
			return
		}
		utils.Logger.Errorf("Error al eliminar perfil %d para usuario %d: %v", profileID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el perfil, detalles: " + err.Error()})
		return
	}

	utils.Logger.Infof("Perfil %d eliminado por usuario %d", profileID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Perfil eliminado"})
}

func (h *ProfileHandler) Update(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warnf("Error al parsear JSON en Update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Warnf("ID de perfil inválido: %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido", "details": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	err = h.Usecase.Update(uint(profileID), userID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Logger.Warnf("Perfil %d no encontrado o acceso denegado para usuario %d", profileID, userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Perfil no encontrado o no permitido"})
			return
		}
		utils.Logger.Errorf("Error al actualizar perfil %d para usuario %d: %v", profileID, userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el perfil", "details": err.Error()})
		return
	}

	utils.Logger.Infof("Perfil %d actualizado por usuario %d", profileID, userID)
	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado"})
}
