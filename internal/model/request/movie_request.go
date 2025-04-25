package request

type CreateMovieRequest struct {
	Title       string
	Description string
	Year        int
	ImageURL    string
	StreamURL   string
	Genre       string
	Duration    int
}

type UpdateMovieRequest struct {
	Title       *string
	Description *string
	Year        *int
	ImageURL    *string
	StreamURL   *string
	Genre       *string
	Duration    *int
}
