package response

type MovieResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        int    `json:"year"`
	ImageURL    string `json:"image_url"`
	StreamURL   string `json:"stream_url"`
	Genre       string `json:"genre"`
	Duration    int    `json:"duration"`
}
