package handler

import (
	"net/http"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type PartyHandler struct {
	PartyUsecase *usecase.PartyUsecase
}

func NewPartyHandler(partyUsecase *usecase.PartyUsecase) *PartyHandler {
	return &PartyHandler{PartyUsecase: partyUsecase}
}

// CreateParty godoc
// @Summary Crear una nueva party
// @Description Crea una nueva party con el host y película especificados
// @Tags Parties
// @Accept json
// @Produce json
// @Param request body request.CreatePartyRequest true "Datos para crear la party"
// @Success 201 {object} response.SuccessResponse "Party creada exitosamente"
// @Failure 400 {object} response.ErrorResponse "Entrada inválida"
// @Failure 500 {object} response.ErrorResponse "Error al crear la party"
// @Router /api/party [post]
func (h *PartyHandler) CreateParty(c *gin.Context) {
	var req request.CreatePartyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorf("Error binding CreatePartyRequest: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	party := &domain.Party{
		HostID:  req.HostID,
		MovieID: req.MovieID,
	}

	err := h.PartyUsecase.CreateParty(c.Request.Context(), party)
	if err != nil {
		utils.Logger.Errorf("Error creating party: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create party"})
		return
	}

	c.JSON(http.StatusCreated, response.CreatePartyResponse{
		Message:   "Party created successfully",
		PartyCode: party.PartyCode,
	})
}

// JoinParty godoc
// @Summary Unirse a una party existente
// @Description Permite a un perfil unirse a una party usando su código
// @Tags Parties
// @Accept json
// @Produce json
// @Param code path string true "Código de la party"
// @Param profile_id header uint true "ID del perfil"
// @Success 200 {object} response.SuccessResponse "Unión exitosa a la party"
// @Failure 400 {object} response.ErrorResponse "Código de party o ID de perfil faltante"
// @Failure 500 {object} response.ErrorResponse "Error al unirse a la party"
// @Router /api/party/join/{code} [post]
func (h *PartyHandler) JoinParty(c *gin.Context) {
	partyCode := c.Param("code")
	profileID := c.GetUint("profile_id")

	if partyCode == "" {
		utils.Logger.Warn("Party code is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Party code is required"})
		return
	}

	if profileID == 0 {
		utils.Logger.Warn("Profile ID is missing in context")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile ID is required"})
		return
	}

	err := h.PartyUsecase.JoinParty(c.Request.Context(), partyCode, profileID)
	if err != nil {
		utils.Logger.Errorf("Error joining party: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join party"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Joined party successfully"})
}
