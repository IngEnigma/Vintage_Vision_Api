package handler

import (
	"net/http"
	"strconv"

	"vintage-vision-api/internal/cloudinary"
	"vintage-vision-api/internal/constants"
	"vintage-vision-api/internal/model/request"
	"vintage-vision-api/internal/model/response"
	"vintage-vision-api/internal/usecase"
	"vintage-vision-api/internal/utils"
	"vintage-vision-api/repository/mapper"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	Usecase *usecase.MovieUsecase
}

func NewMovieHandler(u *usecase.MovieUsecase) *MovieHandler {
	return &MovieHandler{Usecase: u}
}

// GetAll godoc
// @Summary Obtener todas las películas
// @Description Obtiene un listado paginado de películas
// @Tags Movies
// @Produce json
// @Param limit query int false "Límite de resultados (default 10)" default(10)
// @Param offset query int false "Desplazamiento (default 0)" default(0)
// @Success 200 {array} response.MovieResponse "Listado de películas"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/movies [get]
func (h *MovieHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()

	movies, err := h.Usecase.GetAll(ctx, 10, 0)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgGetMovies, err)
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

// GetMoviesByGenre godoc
// @Summary Obtener películas por género
// @Description Obtiene películas filtradas por género con paginación
// @Tags Movies
// @Produce json
// @Param genre query string true "Género para filtrar"
// @Param page query int false "Página (default 1)" default(1)
// @Param limit query int false "Límite por página (default 10)" default(10)
// @Success 200 {object} object{message=string,genre=string,page=int,limit=int,data=[]response.MovieResponse} "Resultados paginados"
// @Failure 400 {object} response.ErrorResponse "Parámetros inválidos"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/movies/by-genre [get]
func (h *MovieHandler) GetMoviesByGenre(c *gin.Context) {
	genre := c.Query("genre")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		utils.Logger.Warnf("%s: %v", constants.ErrMsgInvalidPage, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrMsgInvalidPage})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.Logger.Warnf("%s: %v", constants.ErrMsgInvalidLimit, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrMsgInvalidLimit})
		return
	}

	movies, err := h.Usecase.GetMoviesByGenre(c.Request.Context(), genre, page, limit)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetMoviesByGenre, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": constants.MsgMoviesRetrievedByGenre,
		"genre":   genre,
		"page":    page,
		"limit":   limit,
		"data":    movies,
	})
}

// GetMovieDetail godoc
// @Summary Obtener detalles de película
// @Description Obtiene información detallada de una película específica
// @Tags Movies
// @Produce json
// @Param id path int true "ID de la película"
// @Success 200 {object} response.MovieDetailResponse "Detalles completos"
// @Failure 400 {object} response.ErrorResponse "ID inválido"
// @Failure 404 {object} response.ErrorResponse "No encontrado"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/movies/{id}/detail [get]
func (h *MovieHandler) GetMovieDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidID, err)
		return
	}

	movie, err := h.Usecase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgGetMovie, err)
		return
	}

	detailResponse := response.MovieDetailResponse{
		ID:          strconv.FormatUint(uint64(movie.ID), 10),
		Description: movie.Description,
		Genre:       movie.Genre,
		Year:        movie.Year,
		ImageURL:    movie.ImageURL,
	}

	c.JSON(http.StatusOK, detailResponse)
}

// GetMoviePlayer godoc
// @Summary Obtener datos para reproductor
// @Description Obtiene la información necesaria para reproducir la película
// @Tags Movies
// @Produce json
// @Param id path int true "ID de la película"
// @Success 200 {object} response.MoviePlayerResponse "Datos para reproductor"
// @Failure 400 {object} response.ErrorResponse "ID inválido"
// @Failure 404 {object} response.ErrorResponse "No encontrado"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/movies/{id}/player [get]
func (h *MovieHandler) GetMoviePlayer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidID, err)
		return
	}

	movie, err := h.Usecase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgGetMovie, err)
		return
	}

	playerResponse := response.MoviePlayerResponse{
		ID:        strconv.FormatUint(uint64(movie.ID), 10),
		Title:     movie.Title,
		StreamURL: movie.StreamURL,
	}

	c.JSON(http.StatusOK, playerResponse)
}

// GetPreviewByGenre godoc
// @Summary Vista previa por género
// @Description Obtiene una vista previa de películas por género
// @Tags Movies
// @Produce json
// @Param genre query string true "Género para filtrar"
// @Param page query int false "Página (default 1)" default(1)
// @Param limit query int false "Límite por página (default 10)" default(10)
// @Success 200 {object} object{message=string,genre=string,page=int,limit=int,data=[]response.MoviePreviewResponse} "Vista previa"
// @Failure 400 {object} response.ErrorResponse "Parámetros inválidos"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/movies/preview [get]
func (h *MovieHandler) GetPreviewByGenre(c *gin.Context) {
	genre := c.Query("genre")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		utils.Logger.Warnf("%s: %v", constants.ErrMsgInvalidPage, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrMsgInvalidPage})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		utils.Logger.Warnf("%s: %v", constants.ErrMsgInvalidLimit, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrMsgInvalidLimit})
		return
	}

	movies, err := h.Usecase.GetMoviesByGenre(c.Request.Context(), genre, page, limit)
	if err != nil {
		utils.Logger.Errorf("%s: %v", constants.ErrMsgGetMoviesByGenre, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	previews := mapper.ToMoviePreviewList(movies)

	c.JSON(http.StatusOK, gin.H{
		"message": constants.MsgMoviesRetrievedByGenre,
		"genre":   genre,
		"page":    page,
		"limit":   limit,
		"data":    previews,
	})
}

// GetByID godoc
// @Summary Obtener película completa
// @Description Obtiene todos los datos de una película
// @Tags Movies
// @Produce json
// @Param id path int true "ID de la película"
// @Success 200 {object} response.MovieResponse "Datos completos"
// @Failure 400 {object} response.ErrorResponse "ID inválido"
// @Failure 404 {object} response.ErrorResponse "No encontrado"
// @Failure 500 {object} response.ErrorResponse "Error interno del servidor"
// @Router /api/movies/{id} [get]
func (h *MovieHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidID, err)
		return
	}

	movie, err := h.Usecase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgGetMovie, err)
		return
	}

	c.JSON(http.StatusOK, movie)
}

// Create godoc
// @Summary Crear película
// @Description Crea una nueva película con sus archivos multimedia
// @Tags Movies
// @Accept mpfd
// @Produce json
// @Param title formData string true "Título"
// @Param description formData string true "Descripción"
// @Param year formData int true "Año"
// @Param genre formData string true "Género"
// @Param duration formData int true "Duración (minutos)"
// @Param image formData file true "Imagen de portada"
// @Param video formData file true "Archivo de video"
// @Success 201 {object} response.SuccessResponse "Película creada"
// @Failure 400 {object} response.ErrorResponse "Datos inválidos"
// @Failure 500 {object} response.ErrorResponse "Error al crear"
// @Router /api/movies [post]
func (h *MovieHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	title := c.PostForm("title")
	description := c.PostForm("description")
	year, err := strconv.Atoi(c.PostForm("year"))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidYear, err)
		return
	}
	genre := c.PostForm("genre")
	duration, err := strconv.Atoi(c.PostForm("duration"))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidDuration, err)
		return
	}

	imageFile, err := c.FormFile("image")
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgUploadImage, err)
		return
	}
	videoFile, err := c.FormFile("video")
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgUploadVideo, err)
		return
	}

	image, err := imageFile.Open()
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgOpenImage, err)
		return
	}
	defer image.Close()

	video, err := videoFile.Open()
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgOpenVideo, err)
		return
	}
	defer video.Close()

	imageURL, err := cloudinary.UploadImageToCloudinary(image, imageFile.Filename)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgUploadImage, err)
		return
	}
	videoURL, err := cloudinary.UploadVideoToCloudinary(video, videoFile.Filename)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgUploadVideo, err)
		return
	}

	req := request.CreateMovieRequest{
		Title:       title,
		Description: description,
		Year:        year,
		ImageURL:    imageURL,
		StreamURL:   videoURL,
		Genre:       genre,
		Duration:    duration,
	}

	if _, err := h.Usecase.Create(ctx, req); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgCreateMovie, err)
		return
	}

	utils.Logger.Infof(constants.MsgMovieCreatedSuccessfully)
	c.JSON(http.StatusCreated, gin.H{"message": constants.MsgMovieCreatedSuccessfully})
}

// Update godoc
// @Summary Actualizar película
// @Description Actualiza los datos de una película existente
// @Tags Movies
// @Accept mpfd
// @Produce json
// @Param id path int true "ID de la película"
// @Param title formData string false "Nuevo título"
// @Param description formData string false "Nueva descripción"
// @Param year formData int false "Nuevo año"
// @Param genre formData string false "Nuevo género"
// @Param duration formData int false "Nueva duración"
// @Param image formData file false "Nueva imagen"
// @Param video formData file false "Nuevo video"
// @Success 200 {object} response.SuccessResponse "Película actualizada"
// @Failure 400 {object} response.ErrorResponse "Datos inválidos"
// @Failure 404 {object} response.ErrorResponse "No encontrado"
// @Failure 500 {object} response.ErrorResponse "Error al actualizar"
// @Router /api/movies/{id} [put]
func (h *MovieHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidID, err)
		return
	}

	if c.Request.MultipartForm == nil {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgParseForm, err)
			return
		}
	}

	var req request.UpdateMovieRequest

	if val := c.PostForm("title"); val != "" {
		req.Title = &val
	}
	if val := c.PostForm("description"); val != "" {
		req.Description = &val
	}
	if val := c.PostForm("year"); val != "" {
		if year, err := strconv.Atoi(val); err == nil {
			req.Year = &year
		} else {
			utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidYear, err)
			return
		}
	}
	if val := c.PostForm("genre"); val != "" {
		req.Genre = &val
	}
	if val := c.PostForm("duration"); val != "" {
		if dur, err := strconv.Atoi(val); err == nil {
			req.Duration = &dur
		} else {
			utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidDuration, err)
			return
		}
	}

	if imageFile, err := c.FormFile("image"); err == nil {
		image, _ := imageFile.Open()
		defer image.Close()

		imageURL, err := cloudinary.UploadImageToCloudinary(image, imageFile.Filename)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgUploadImage, err)
			return
		}
		req.ImageURL = &imageURL
	}

	if videoFile, err := c.FormFile("video"); err == nil {
		video, _ := videoFile.Open()
		defer video.Close()

		videoURL, err := cloudinary.UploadVideoToCloudinary(video, videoFile.Filename)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgUploadVideo, err)
			return
		}
		req.StreamURL = &videoURL
	}

	if _, err := h.Usecase.Update(ctx, uint(id), req); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgUpdateMovie, err)
		return
	}

	utils.Logger.Infof(constants.MsgMovieUpdatedSuccessfully)
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgMovieUpdatedSuccessfully})
}

// Delete godoc
// @Summary Eliminar película
// @Description Elimina una película existente
// @Tags Movies
// @Produce json
// @Param id path int true "ID de la película"
// @Success 200 {object} response.SuccessResponse "Película eliminada"
// @Failure 400 {object} response.ErrorResponse "ID inválido"
// @Failure 404 {object} response.ErrorResponse "No encontrado"
// @Failure 500 {object} response.ErrorResponse "Error al eliminar"
// @Router /api/movies/{id} [delete]
func (h *MovieHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, constants.ErrMsgInvalidID, err)
		return
	}

	err = h.Usecase.Delete(ctx, uint(id))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, constants.ErrMsgDeleteMovie, err)
		return
	}

	utils.Logger.Infof(constants.MsgMovieDeletedSuccessfully)
	c.JSON(http.StatusOK, gin.H{"message": constants.MsgMovieDeletedSuccessfully})
}
