package response

// MovieResponse representa la estructura de respuesta para una película.
// swagger:model MovieResponse
type MovieResponse struct {
	// ID único de la película
	// required: true
	ID uint `json:"id"`

	// Título de la película
	// required: true
	Title string `json:"title"`

	// Descripción de la película
	// required: true
	Description string `json:"description"`

	// Año de lanzamiento
	// required: true
	Year int `json:"year"`

	// URL de la imagen de portada
	// required: true
	ImageURL string `json:"image_url"`

	// URL del video para streaming
	// required: true
	StreamURL string `json:"stream_url"`

	// Género de la película
	// required: true
	Genre string `json:"genre"`

	// Duración en minutos
	// required: true
	Duration int `json:"duration"`
}
