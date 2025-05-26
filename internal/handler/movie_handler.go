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

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	Usecase *usecase.MovieUsecase
}

func NewMovieHandler(u *usecase.MovieUsecase) *MovieHandler {
	return &MovieHandler{Usecase: u}
}

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
