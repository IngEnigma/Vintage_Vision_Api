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

func NewAuthHandler(u *usecase.UserUsecase) *AuthHandler {
	return &AuthHandler{Usecase: u}
}

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
