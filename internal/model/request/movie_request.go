package request

type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Year        int    `json:"year" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required,url"`
	StreamURL   string `json:"stream_url" binding:"required,url"`
	Genre       string `json:"genre" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
}

type UpdateMovieRequest struct {
	Title       *string `json:"title" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
	Year        *int    `json:"year" binding:"omitempty"`
	ImageURL    *string `json:"image_url" binding:"omitempty,url"`
	StreamURL   *string `json:"stream_url" binding:"omitempty,url"`
	Genre       *string `json:"genre" binding:"omitempty"`
	Duration    *int    `json:"duration" binding:"omitempty"`
}
