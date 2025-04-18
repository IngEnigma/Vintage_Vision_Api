package handler

import (
	"net/http"
	"strconv"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	Usecase *usecase.MovieUsecase
}

func NewMovieHandler(u *usecase.MovieUsecase) *MovieHandler {
	return &MovieHandler{Usecase: u}
}

func (h *MovieHandler) GetAll(c *gin.Context) {
	movies, err := h.Usecase.GetAll()
	if err != nil {
		utils.Logger.Errorf("Error al obtener las películas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las películas, detalles: " + err.Error()})
		return
	}

	var res []response.MovieResponse
	for _, m := range movies {
		res = append(res, response.MovieResponse{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Year:        m.Year,
			ImageURL:    m.ImageURL,
			StreamURL:   m.StreamURL,
			Genre:       m.Genre,
			Duration:    m.Duration,
		})
	}

	utils.Logger.Infof("Películas obtenidas con éxito")
	c.JSON(http.StatusOK, res)
}

func (h *MovieHandler) Create(c *gin.Context) {
	var req request.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warnf("Error al bindear los datos para crear la película: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos, detalles: " + err.Error()})
		return
	}

	if err := h.Usecase.Create(req); err != nil {
		utils.Logger.Errorf("Error al crear la película: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la película, detalles: " + err.Error()})
		return
	}

	utils.Logger.Infof("Película creada: %s", req.Title)
	c.JSON(http.StatusCreated, gin.H{"message": "Película creada"})
}

func (h *MovieHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Logger.Warnf("ID inválido para actualizar película: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, detalles: " + err.Error()})
		return
	}

	var req request.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Warnf("Error al bindear los datos para actualizar la película: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos, detalles: " + err.Error()})
		return
	}

	if err := h.Usecase.Update(uint(id), req); err != nil {
		utils.Logger.Errorf("Error al actualizar la película con ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la película, detalles: " + err.Error()})
		return
	}

	utils.Logger.Infof("Película con ID %d actualizada con éxito", id)
	c.JSON(http.StatusOK, gin.H{"message": "Película actualizada"})
}

func (h *MovieHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.Logger.Warnf("ID inválido para eliminar película: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido, detalles: " + err.Error()})
		return
	}

	err = h.Usecase.Delete(uint(id))
	if err != nil {
		utils.Logger.Errorf("Error al eliminar la película con ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la película, detalles: " + err.Error()})
		return
	}

	utils.Logger.Infof("Película con ID %d eliminada correctamente", id)
	c.JSON(http.StatusOK, gin.H{"message": "Película eliminada correctamente"})
}
