package handler

import (
	"errors"
	"net/http"
	"strconv"

	"vintage-vision-api/internal/constants"
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

// NewProfileHandler crea una nueva instancia de ProfileHandler
func NewProfileHandler(u *usecase.ProfileUsecase) *ProfileHandler {
	return &ProfileHandler{Usecase: u}
}

// Create godoc
// @Summary Crea un nuevo perfil
// @Description Crea un perfil asociado al usuario autenticado
// @Tags profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body request.CreateProfileRequest true "Datos del perfil a crear"
// @Success 201 {object} map[string]string "Perfil creado exitosamente"
// @Router /api/profiles [post]
func (h *ProfileHandler) Create(c *gin.Context) {
	var req request.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidRequest, err)
		return
	}

	userID := c.GetUint("user_id")
	if err := h.Usecase.Create(c.Request.Context(), userID, req); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgRegisterProfile, err)
		return
	}

	utils.Logger.Infof(constants.MsgProfileCreatedSuccessfully+" UserID: %d", userID)
	c.JSON(http.StatusCreated, gin.H{"message": constants.MsgProfileCreatedSuccessfully})
}

// GetAll godoc
// @Summary Lista todos los perfiles del usuario autenticado
// @Description Obtiene todos los perfiles asociados a un usuario
// @Tags profiles
// @Produce json
// @Security BearerAuth
// @Success 200 {array} response.ProfileResponse "Lista de perfiles"
// @Router /api/profiles [get]
func (h *ProfileHandler) GetAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	profiles, err := h.Usecase.GetAll(c.Request.Context(), userID)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgGetProfiles, err)
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

	utils.Logger.Infof(constants.MsgProfilesRetrievedSuccessfully+" UserID: %d", userID)
	c.JSON(http.StatusOK, res)
}

// Delete godoc
// @Summary Elimina un perfil
// @Description Elimina un perfil por ID asociado al usuario autenticado
// @Tags profiles
// @Security BearerAuth
// @Param id path int true "ID del perfil a eliminar"
// @Success 200 {object} map[string]string "Perfil eliminado exitosamente"
// @Router /api/profiles/{id} [delete]
func (h *ProfileHandler) Delete(c *gin.Context) {
	profileIDStr := c.Param("id")
	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidID, err)
		return
	}

	userID := c.GetUint("user_id")
	err = h.Usecase.Delete(c.Request.Context(), uint(profileID), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, http.StatusNotFound, constants.ErrMsgProfileNotFound, err)
			return
		}
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgDeleteProfile, err)
		return
	}

	utils.Logger.Infof(constants.MsgProfileDeletedSuccessfully+" ProfileID: %d, UserID: %d", profileID, userID)
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgProfileDeletedSuccessfully})
}

// Update godoc
// @Summary Actualiza un perfil
// @Description Actualiza los datos de un perfil específico
// @Tags profiles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID del perfil a actualizar"
// @Param profile body request.UpdateProfileRequest true "Datos a actualizar"
// @Success 200 {object} map[string]string "Perfil actualizado exitosamente"
// @Failure 400 {object} response.ErrorResponse "Solicitud inválida o ID incorrecto"
// @Router /api/profiles/{id} [patch]
func (h *ProfileHandler) Update(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidRequest, err)
		return
	}

	profileID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidID, err)
		return
	}

	userID := c.GetUint("user_id")
	err = h.Usecase.Update(c.Request.Context(), uint(profileID), userID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.HandleError(c, http.StatusNotFound, constants.ErrMsgProfileNotFound, err)
			return
		}
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgUpdateProfile, err)
		return
	}

	utils.Logger.Infof(constants.MsgProfileUpdatedSuccessfully+" ProfileID: %d, UserID: %d", profileID, userID)
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgProfileUpdatedSuccessfully})
}
