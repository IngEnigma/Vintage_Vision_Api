package handler

import (
	"net/http"
	"vintage-vision-api/internal/domain"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type PartyHandler struct {
	PartyUsecase usecase.PartyUsecase
}

func NewPartyHandler(partyUsecase usecase.PartyUsecase) *PartyHandler {
	return &PartyHandler{
		PartyUsecase: partyUsecase,
	}
}

func (h *PartyHandler) CreateParty(c *gin.Context) {
	var party domain.Party
	if err := c.ShouldBindJSON(&party); err != nil {
		utils.Logger.Errorf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if party.HostID == 0 || party.MovieID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "HostID and MovieID are required"})
		return
	}

	err := h.PartyUsecase.CreateParty(c.Request.Context(), &party)
	if err != nil {
		utils.Logger.Errorf("Error creating party: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create party"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Party created successfully", "party_code": party.PartyCode})
}

func (h *PartyHandler) JoinParty(c *gin.Context) {
	partyCode := c.Param("code")
	profileID := c.GetUint("profile_id")

	if partyCode == "" {
		utils.Logger.Warn("Party code is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Party code is required"})
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
