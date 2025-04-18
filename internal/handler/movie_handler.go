package handler

import (
	"net/http"
	"strconv"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"

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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener las películas"})
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

	c.JSON(http.StatusOK, res)
}

func (h *MovieHandler) Create(c *gin.Context) {
	var req request.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.Usecase.Create(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la película"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Película creada"})
}

func (h *MovieHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req request.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	if err := h.Usecase.Update(uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar la película"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Película actualizada"})
}

func (h *MovieHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = h.Usecase.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la película"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Película eliminada correctamente"})
}
